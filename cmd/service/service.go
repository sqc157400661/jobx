package service

import (
	"fmt"
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/pkg/mysql"
	"github.com/sqc157400661/jobx/pkg/options/flowopt"
	"sync"
	"time"

	"k8s.io/klog/v2"

	"github.com/sqc157400661/jobx/cmd/log"
	"github.com/sqc157400661/jobx/internal"
	"github.com/sqc157400661/jobx/internal/collector"
	"github.com/sqc157400661/jobx/internal/queue"
	joberrors "github.com/sqc157400661/jobx/pkg/errors"
	"github.com/sqc157400661/jobx/pkg/model"
	"github.com/sqc157400661/jobx/pkg/providers"
)

const (
	processJobEnQueueInterval = 1 * time.Second
	localQueueLen             = 50
	collectorJobLen           = 10
)

// JobFlow represents a job processing workflow that steals jobs from database,
// manages local queue, and coordinates job distribution through tracker
type JobFlow struct {
	// Unique identifier for the JobFlow instance
	uid string
	// Configuration options
	opts flowopt.Options
	// Job collector from database
	collector collector.Collector
	// Local buffer for stolen jobs
	localQueue *queue.TaskQueue
	// Worker pool for job execution
	worker internal.WorkerPool
	// Tracks job progress and status
	tracker internal.Tracker
	// cronTrigger job cron trigger
	cronTrigger collector.CronTrigger

	// Concurrency control
	stopChan  chan struct{}
	startOnce sync.Once
	stopOnce  sync.Once
	stopped   chan struct{}
}

// NewJobFlow creates a new JobFlow instance with specified UID and database connection
func NewJobFlow(uid string, conf config.MySQL, opts ...flowopt.OptionFunc) (jf *JobFlow, err error) {
	if uid == "" || conf.Host == "" {
		err = joberrors.ErrInvalidParameter
		return
	}
	err = mysql.SetDB(conf)
	if err != nil {
		return nil, err
	}
	jf = &JobFlow{
		uid:        uid,
		localQueue: queue.NewTaskQueue(localQueueLen),
		stopChan:   make(chan struct{}),
		stopped:    make(chan struct{}),
	}
	jf.opts = flowopt.DefaultOption
	for _, opt := range opts {
		opt(&jf.opts)
	}
	// Initialize components with proper isolation
	jf.collector = collector.NewDefaultCollector(uid, jf.opts.Tenant, jf.opts.AppName, collectorJobLen)
	jf.worker = internal.NewDefaultWorkerPool(jf.opts.PoolLen)
	jf.tracker = internal.NewTracker(jf.worker, jf.localQueue)
	if !jf.opts.DisableCron {
		jf.cronTrigger = collector.NewDefaultCronTrigger()
	}
	return
}

// Register adds task providers to the global provider registry
func (jf *JobFlow) Register(p ...providers.TaskProvider) (err error) {
	for _, v := range p {
		if v != nil {
			name := v.Name()
			if providers.Has(name) {
				err = fmt.Errorf("the taskProvider already exists with name %s", name)
			}
			providers.Set(v)
		}
	}
	return
}

// AddProvider registers a task provider with optional alias
func (jf *JobFlow) AddProvider(t providers.TaskProvider, action ...string) (err error) {
	if t == nil {
		return
	}
	name := t.Name()
	if len(action) > 0 {
		name = action[0]
	}
	if providers.Has(name) {
		err = fmt.Errorf("the taskProvider already exists with name %s", name)
	}
	providers.Set(t)
	return
}

// Start initiates the job processing workflow
func (jf *JobFlow) Start() {
	jf.startOnce.Do(func() {
		go func() {
			jf.tracker.Start()
			if jf.cronTrigger != nil {
				jf.cronTrigger.Start()
			}
			go jf.processJob()
			klog.Infof("jobFlow:%s ,uid:%s,starting", jf.opts.Desc, jf.uid)
		}()
	})
}

// processJobEnqueue manages job stealing and queue population
func (jf *JobFlow) processJob() {
	var ticker *time.Ticker
	ticker = time.NewTicker(processJobEnQueueInterval)
	defer ticker.Stop()
	for {
		select {
		case <-jf.tracker.Waiting():
			jf.stealJobs()
		case <-ticker.C:
			jf.stealJobs()
			jf.stealCronJobs()
		case <-jf.stopChan:
			err := jf.collector.ReleaseJobs()
			if err != nil {
				klog.Error(err)
			}
			close(jf.stopped)
			return
		}
	}
}

func (jf *JobFlow) RunSyncJob(Job *model.Job) (err error) {
	return jf.tracker.StartJob(Job, true)
}

// stealJobs retrieves jobs from collector and enqueues them
func (jf *JobFlow) stealJobs() {
	if jf.localQueue.PendingCount() > localQueueLen/2 {
		return
	}
	jobs, err := jf.collector.StealJobs()
	if err != nil {
		klog.Errorf("jobFlow:%s ,uid:%s,StealJobs job Err:%s", jf.opts.Desc, jf.uid, err.Error())
		return
	}
	err = jf.EnqueueJobs(jobs...)
	if err != nil {
		klog.Errorf("jobFlow:%s ,uid:%s,EnqueueJobs job Err:%s", jf.opts.Desc, jf.uid, err.Error())
		return
	}
}

// stealCronJobs retrieves cron jobs from collector
func (jf *JobFlow) stealCronJobs() {
	if jf.cronTrigger == nil {
		return
	}
	cronJobs, err := jf.collector.StealCronJobs()
	if err != nil {
		klog.Errorf("jobFlow:%s ,uid:%s,StealCronJobs job Err:%s", jf.opts.Desc, jf.uid, err.Error())
		return
	}
	jf.cronTrigger.Load(cronJobs)
}

// EnqueueJobs adds jobs to the local queue with overflow handling
func (jf *JobFlow) EnqueueJobs(jobs ...*model.Job) error {
	var errs []error
	for _, job := range jobs {
		err := jf.localQueue.AddToFront(job)
		if err != nil {
			errs = append(errs, err)
			if joberrors.IsQueueFull(err) {
				releaseErr := jf.collector.ReleaseJobByID(job.ID)
				if releaseErr != nil {
					klog.Error(releaseErr)
				}
			} else {
				klog.Error(err)
				return err
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("enqueue errors encountered: %v", errs)
	}
	return nil
}

// Stop gracefully shuts down the JobFlow
func (jf *JobFlow) Stop() {
	jf.stopOnce.Do(func() {
		if jf.cronTrigger != nil {
			jf.cronTrigger.Stop()
		}
		jf.tracker.Stop()
		close(jf.stopChan)
		<-jf.stopped
	})
	// todo 逻辑优化 退出时log刷盘暂时放到这里
	if log.Logger != nil {
		if err := log.Logger.Flush(); err != nil {
			klog.Errorf("Log flush failed: %v", err)
		}
	}
}

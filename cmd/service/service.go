package service

import (
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/sqc157400661/jobx/cmd/log"
	"github.com/sqc157400661/jobx/internal"
	"github.com/sqc157400661/jobx/internal/collector"
	"github.com/sqc157400661/jobx/pkg/dao"
	"github.com/sqc157400661/jobx/pkg/options"
	"github.com/sqc157400661/jobx/pkg/providers"
	"k8s.io/klog/v2"
	"sync"
	"time"
)

/*
	todo 优化，
	1.同一台机器上运行多个uid相同的JobFlow会有问题
	2.任务超时 释放锁机制
	3.优化log，记录log
*/

type JobFlow struct {
	// Unique ID for JobFlow
	uid       string
	opts      options.Options
	collector collector.Collector
	worker    internal.Worker
	tracker   internal.Tracker
	// worker stop signal
	stopChan chan struct{}
	// The start operation is performed only once
	startOnce sync.Once
	// The stop operation is performed only once
	stopOnce sync.Once
	stopped  chan struct{}
}

func NewJobFlow(uid string, db *xorm.Engine, opts ...options.OptionFunc) (jf *JobFlow, err error) {
	if uid == "" || db == nil {
		err = errors.New("NewJobFlow err,param is not available")
		return
	}
	jf = &JobFlow{
		uid:      uid,
		stopChan: make(chan struct{}),
		stopped:  make(chan struct{}),
	}
	o := options.DefaultOption
	for _, opt := range opts {
		opt(&o)
	}
	jf.opts = o
	// task execution occupies a separate session connection
	dao.JFDb = db
	jf.collector = collector.NewDefaultCollector(db, uid, jf.opts.PoolLen)
	jf.worker = internal.NewDefaultWorker(o.PoolLen)
	jf.tracker = internal.NewTracker(jf.worker)
	return
}

// 注册action
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

func (jf *JobFlow) Start() {
	var ticker *time.Ticker
	jf.startOnce.Do(func() {
		go func() {
			jf.worker.Run()
			klog.Infof("jobFlow:%s ,uid:%s,starting", jf.opts.Desc, jf.uid)
			ticker = time.NewTicker(jf.opts.LoopInterval)
			for {
				select {
				case <-ticker.C:
					// 这里去尝试获取资源并执行锁定
					jobs, err := jf.collector.StealJob()
					if err != nil {
						klog.Errorf("jobFlow:%s ,uid:%s,collector job Err:%s", jf.opts.Desc, jf.uid, err.Error())
					}
					for _, job := range jobs {
						err = jf.tracker.AddRootJob(job, jf.uid)
						if err != nil {
							klog.Error(err)
						}
					}
				case <-jf.stopChan:
					err := jf.collector.ReleaseJob()
					if err != nil {
						klog.Error(err)
					}
					close(jf.stopped)
					return
				}
			}
		}()
	})
}

func (jf *JobFlow) Quit() {
	jf.stopOnce.Do(func() {
		jf.worker.Quit()
		close(jf.stopChan)
		<-jf.stopped
	})
	// todo 逻辑优化 退出时log刷盘暂时放到这里
	if log.Logger != nil {
		_ = log.Logger.Flush()
	}
}

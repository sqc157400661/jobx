package internal

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"k8s.io/klog/v2"

	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/internal/queue"
	"github.com/sqc157400661/jobx/pkg/dao"
	joberrors "github.com/sqc157400661/jobx/pkg/errors"
)

type Tracker interface {
	// Start begins processing jobs from the queue
	Start()
	// Stop terminates the tracker and releases resources
	Stop()
	// Waiting returns a channel that signals when the queue is empty
	Waiting() <-chan struct{}
	// StartJob initiates tracking and execution of a root job
	StartJob(rootJob *dao.Job, isSync bool) (err error)
}

type DefaultTracker struct {
	worker WorkerPool
	// Map of root job IDs to their tracking items
	trackerMap map[int]*trackItem
	// Local job queue for processing
	localQueue *queue.TaskQueue
	// Signals when queue becomes empty
	waitChan chan struct{}
	// Signals tracker shutdown
	stopChan chan struct{}
	// Synchronizes access to trackerMap
	mutex sync.RWMutex
	// Channel for receiving job status updates
	TrackSignal chan TrackSignal
}

type TrackSignal struct {
	RootID      int    // Root job ID
	JobID       int    // Specific job ID reporting status
	IsSucceeded bool   // True if job succeeded
	IsPaused    bool   // True if job is paused
	Msg         string // Additional status message
}

// NewTracker creates and initializes a new DefaultTracker instance
// todo :When something goes wrong here, we need to capture it and alert it
func NewTracker(worker WorkerPool, queue *queue.TaskQueue) (t *DefaultTracker) {
	t = &DefaultTracker{
		stopChan:    make(chan struct{}),
		worker:      worker,
		localQueue:  queue,
		trackerMap:  make(map[int]*trackItem),
		waitChan:    make(chan struct{}),
		TrackSignal: make(chan TrackSignal, 100),
	}
	go t.track()
	return
}

func (t *DefaultTracker) Start() {
	t.worker.Start()
	go func() {
		for {
			select {
			case <-t.stopChan:
				return
			default:
				job, err := t.localQueue.Dequeue()
				if joberrors.IsQueueEmpty(err) {
					t.waitChan <- struct{}{}
					time.Sleep(1 * time.Second)
					continue
				}
				if err != nil {
					klog.Error("Dequeue error:", err)
					continue
				}
				if err = t.StartJob(job, false); err != nil {
					klog.Error("StartJob error:", err)
				}
			}

		}
	}()
}

// Waiting returns the wait status channel
func (t *DefaultTracker) Waiting() <-chan struct{} {
	return t.waitChan
}

// StartJob initializes tracking for a new root job
func (t *DefaultTracker) StartJob(rootJob *dao.Job, isSync bool) (err error) {
	// if it already exists in tracker,return
	if t.get(rootJob.ID) != nil {
		return
	}
	var jobs []dao.Job
	if rootJob.Runnable == config.RunnableYes {
		rootJob.RootID = rootJob.ID
		jobs = []dao.Job{*rootJob}
	} else {
		jobs, err = dao.GetChildRunableJobsByRootId(rootJob.ID)
	}
	if err != nil {
		err = errors.Wrapf(err, "GetChildRunableJobsByRootId err rootID:%d", rootJob.ID)
		return
	}
	if len(jobs) == 0 {
		err = errors.Wrapf(err, "not found RunableJobs rootID:%d", rootJob.ID)
		return
	}
	err = rootJob.MarkRunning()
	if err != nil {
		err = errors.Wrapf(err, "MarkRunning err jobID:%d", rootJob.ID)
		return
	}
	return t.addTrackItem(rootJob.ID, jobs, isSync)
}

func (t *DefaultTracker) addTrackItem(rootId int, jobs []dao.Job, isSync bool) (err error) {
	var ti *trackItem
	ti, err = NewTrackItem(jobs, t.TrackSignal)
	if err != nil {
		return
	}
	if ti.total == ti.finished {
		return dao.UpdateJobStateByID(rootId, &dao.State{
			Phase:  config.PhaseTerminated,
			Status: config.StatusSuccess,
		})
	} else if ti.total < ti.finished {
		err = errors.New(fmt.Sprintf("unknown job status rootID:%d", rootId))
		return
	}
	// Submit jobs to worker pool
	for _, v := range ti.Pipes {
		err = t.worker.Submit(v, isSync)
		if err != nil {
			return
		}
	}
	t.add(rootId, ti)
	return
}

func (t *DefaultTracker) add(key int, v *trackItem) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.trackerMap[key] = v
}
func (t *DefaultTracker) remove(key int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	delete(t.trackerMap, key)
}
func (t *DefaultTracker) get(key int) *trackItem {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.trackerMap[key]
}

// track processes status updates and manages job lifecycle
// if successful, check the entire job chain
func (t *DefaultTracker) track() {
	for trackSignal := range t.TrackSignal {
		rootID := trackSignal.RootID
		jobID := trackSignal.JobID
		msg := trackSignal.Msg
		ti := t.get(rootID)
		if ti == nil {
			// 无主job
			klog.Error("not found trackItem rootID: ", rootID)
			return
		}
		if trackSignal.IsSucceeded {
			ti.DoneOne(jobID)
		} else {
			if trackSignal.IsPaused {
				ti.AddPaused()
			} else {
				ti.FailOne(jobID, errors.New(msg))
			}
		}
		// the life cycle of the root job is end
		// note that only the root job and runnable job have a state, while the job in the middle layer has no state
		if ti.IsFinished() {
			t.localQueue.CompleteTask(rootID)
			var state *dao.State
			var err error
			if ti.IsSucceeded() {
				state = &dao.State{
					Phase:  config.PhaseTerminated,
					Status: config.StatusSuccess,
				}
				err = dao.ReleaseTokens(rootID)
				if err != nil {
					klog.Errorf("ReleaseTokens err:%+v rootID:%d", err, rootID)
				}
			} else if ti.IsPaused() {
				state = &dao.State{
					Phase:  config.PhaseRunning,
					Status: config.StatusPause,
				}
			} else {
				state = &dao.State{
					Phase:  config.PhaseTerminated,
					Status: config.StatusFail,
				}
			}
			t.remove(rootID)
			err = dao.UpdateJobStateByID(rootID, state)
			if err != nil {
				klog.Errorf("TriggerChord err:%+v rootID:%d", err, rootID)
				return
			}
		}
	}
}

// Stop shuts down the tracker
func (t *DefaultTracker) Stop() {
	close(t.stopChan)
	t.worker.Stop()
	close(t.TrackSignal)
}

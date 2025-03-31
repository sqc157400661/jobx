package internal

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"k8s.io/klog/v2"

	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/internal/queue"
	"github.com/sqc157400661/jobx/pkg/dao"
)

/*
负责状态的更新
回调
*/
type Tracker interface {
	AddRootJob(rootJob *dao.Job, uid string) (err error)
}

type DefaultTracker struct {
	worker      Worker
	trackerMap  map[int]*trackItem
	startOnce   sync.Once
	stopChan    chan struct{}
	mutex       sync.RWMutex
	TrackSignal chan TrackSignal
}

type TrackSignal struct {
	RootID      int
	JobID       int
	IsSucceeded bool
	IsPaused    bool
	Msg         string
}

func NewTracker(worker Worker, queue queue.TaskQueue) (t *DefaultTracker) {
	t = &DefaultTracker{
		stopChan:    make(chan struct{}),
		worker:      worker,
		trackerMap:  make(map[int]*trackItem),
		TrackSignal: make(chan TrackSignal, 100),
	}
	go t.track()
	return
}

// Add  add rootJob to Tracker
func (t *DefaultTracker) AddRootJob(rootJob *dao.Job) (err error) {
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
	return t.addTrackItem(rootJob.ID, jobs)
}

func (t *DefaultTracker) addTrackItem(rootId int, jobs []dao.Job) (err error) {
	// todo 这里jobQueue 可以使用对象池来做优化，减少内存分配和回收
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
	for _, v := range ti.Queue {
		err = t.worker.Submit(v)
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

// track listening pipeline results
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

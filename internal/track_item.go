package internal

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/pkg/dao"
	"k8s.io/klog/v2"
	"sync"
)

type trackItem struct {
	Queue    []*Pipeline
	total    int
	failed   int
	paused   int
	finished int
	mutex    sync.RWMutex
}

func NewTrackItem(jobs []dao.Job, trackSignal chan TrackSignal) (ti *trackItem, err error) {
	num := len(jobs)
	if num == 0 {
		return
	}
	ti = &trackItem{
		Queue: make([]*Pipeline, 0, num),
		total: num,
	}
	// 循环处理
	for _, job := range jobs {
		if job.IsReady() {
			err = job.MarkRunning()
			if err != nil {
				ti.FailOne(job.ID, errors.Wrapf(err, "MarkRunning err rootID:%d  jobID:%d", job.RootID, job.ID))
				continue
			}
			var p *Pipeline
			p, err = NewPipeline(job, trackSignal)
			if err != nil {
				return
			}
			ti.Queue = append(ti.Queue, p)
		} else if job.IsFinished() {
			ti.AddSucceeded()
		} else if job.IsPausing() {
			ti.AddPaused()
		} else {
			ti.FailOne(job.ID, errors.New(fmt.Sprintf("unkown status rootID:%d jobID:%d", job.RootID, job.ID)))
		}
	}
	return
}

func (t *trackItem) DoneOne(jobID int) {
	// update job status to Success
	errDb := dao.UpdateJobStateByID(jobID, &dao.State{
		Phase:  config.PhaseTerminated,
		Status: config.StatusSuccess,
	})
	if errDb != nil {
		klog.Error(errDb)
		return
	}
	t.AddSucceeded()
}
func (t *trackItem) AddSucceeded() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.finished += 1
}
func (t *trackItem) AddFailed() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.failed += 1
}
func (t *trackItem) AddPaused() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.paused += 1
}

func (t *trackItem) FailOne(jobID int, err error) {
	// update job status to Fail
	errDb := dao.UpdateJobStateByID(jobID, &dao.State{
		Phase:  config.PhaseTerminated,
		Status: config.StatusFail,
		Reason: err.Error(),
	})
	if errDb != nil {
		klog.Error(errDb)
		return
	}
	t.AddFailed()
}

func (t *trackItem) IsSucceeded() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.finished == t.total
}

func (t *trackItem) IsFinished() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.finished+t.failed+t.paused == t.total
}

func (t *trackItem) IsPaused() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.finished+t.paused == t.total
}

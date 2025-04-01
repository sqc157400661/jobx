package internal

import (
	"sync"

	"k8s.io/klog/v2"

	joberrors "github.com/sqc157400661/jobx/pkg/errors"
	"github.com/sqc157400661/jobx/pkg/model"
)

type Pipeline struct {
	RootID       int
	JobID        int
	Steps        []*model.PipelineTask
	resSignal    chan<- TrackSignal
	mutex        sync.RWMutex
	isPaused     bool
	TotalTask    int // 一共多少个节点任务
	FinishedTask int // 目前已完成的节点个数
}

// NewPipeline 初始化流水线,创建一个group的任务
func NewPipeline(job model.Job, res chan<- TrackSignal) (p *Pipeline, err error) {
	var tasks, readyTasks []*model.PipelineTask
	tasks, err = model.GetPipelineTasksByJobId(job.ID)
	if err != nil {
		return
	}
	if len(tasks) == 0 {
		err = joberrors.NewNoTaskFoundErrorWithJobID(job.RootID, job.ID)
		return
	}
	p = &Pipeline{
		RootID:    job.RootID,
		JobID:     job.ID,
		TotalTask: len(tasks),
		resSignal: res,
	}
	for _, task := range tasks {
		if task.State.IsReady() {
			readyTasks = append(readyTasks, task)
		} else if task.State.IsFinished() {
			p.FinishedTask++
		} else {
			// if the task is pausing or failed, it is considered unable to execute
			// when obtaining a job, it should be verified. so considered that an error has occurred
			p = nil
			err = joberrors.NewNoRunnablePipelineTaskErrorWithJobID(job.RootID, job.ID)
			return nil, err
		}
	}
	p.Steps = readyTasks
	return
}

// check Job status
func (p *Pipeline) IsRunning() bool {
	job, _, err := model.GetJobById(p.JobID)
	if err != nil {
		klog.Error(err)
		return false
	}
	if job.IsPausing() {
		p.Paused()
		return false
	}
	return job.IsRunning()
}

func (p *Pipeline) IsPaused() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.isPaused
}

func (p *Pipeline) Paused() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.isPaused = true
}

func (p *Pipeline) AddSucceeded() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.FinishedTask += 1
}

func (p *Pipeline) Finish() {
	if p.FinishedTask == p.TotalTask || len(p.Steps) == 0 {
		p.writeBack(true, "")
		return
	}
	p.writeBack(false, "pipeline finish err")
}

func (p *Pipeline) Fail(e error) {
	msg := "pipeline fail"
	if e != nil {
		msg = e.Error()
	}
	p.writeBack(false, msg)
	return
}

func (p *Pipeline) writeBack(isSucceeded bool, msg string) {
	p.resSignal <- TrackSignal{
		RootID:      p.RootID,
		JobID:       p.JobID,
		IsSucceeded: isSucceeded,
		IsPaused:    p.isPaused,
		Msg:         msg,
	}
}

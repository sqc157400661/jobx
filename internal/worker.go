package internal

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/sqc157400661/jobx/internal/helper"
	"github.com/sqc157400661/jobx/pkg/dao"
	"github.com/sqc157400661/jobx/pkg/options"
	"github.com/sqc157400661/jobx/pkg/providers"
)

type WorkerPool interface {
	Submit(p *Pipeline) (err error)
	Run()
	Quit()
}

type DefaultWorkerPool struct {
	maxWorkers int
	pipePool   chan *Pipeline
	// The start operation is performed only once
	startOnce sync.Once
	// The stop operation is performed only once
	stopOnce sync.Once
	wg       sync.WaitGroup
	isQuit   bool
}

func NewDefaultWorkerPool(maxWorkers int) (w *DefaultWorkerPool) {
	w = &DefaultWorkerPool{
		maxWorkers: maxWorkers,
		pipePool:   make(chan *Pipeline, maxWorkers*2),
	}
	return
}
func (w *DefaultWorkerPool) Submit(p *Pipeline) (err error) {
	if w.isQuit {
		return errors.New("worker is quit")
	}
	w.pipePool <- p
	return
}

func (w *DefaultWorkerPool) Run() {
	w.startOnce.Do(func() {
		for i := 0; i < w.maxWorkers; i++ {
			w.wg.Add(1)
			go w.worker()
		}
	})
}

func (w *DefaultWorkerPool) worker() {
	defer w.wg.Done()
	for task := range w.pipePool {
		w.process(task)
	}
}

func (w *DefaultWorkerPool) process(pipeline *Pipeline) {
	if pipeline == nil {
		return
	}
	if len(pipeline.Steps) == 0 {
		pipeline.Finish()
		return
	}
	var err error
	var curTask *dao.PipelineTask
	var taskProvider providers.TaskProvider
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("internal error: %v", r)
			fmt.Printf("panic:%s \n", err.Error())
			curTask.Fail(err)
			pipeline.Fail(err)
			w.rollback(taskProvider, curTask.Status)
		}
	}()
	defer func() {
		if err != nil {
			curTask.Fail(err)
			pipeline.Fail(err)
			w.rollback(taskProvider, curTask.Status)
		} else {
			if curTask.IsPausing() {
				pipeline.Paused()
			}
			pipeline.Finish()
		}
	}()
	// TODO : 定期维护执行心跳，用于无主任务回收
	// pipeCtx may interrupt loss
	var pipeCtx map[string]interface{}
	for _, task := range pipeline.Steps {
		curTask = task
		if !pipeline.IsRunning() {
			return
		}
		// find the taskProvider
		taskProvider = providers.Get(curTask.Action)
		if taskProvider == nil {
			err = errors.New(fmt.Sprintf("not found taskProvider:%s", curTask.Action))
			return
		}
		pipeCtx = helper.UnsafeMergeMap(pipeCtx, curTask.Context)
		// 传入相关参数
		curTask.Input = helper.UnsafeMergeMap(curTask.Input, pipeCtx)
		curTask.Context = pipeCtx
		err = curTask.Update()
		if err != nil {
			err = errors.Wrap(err, "write input err")
			return
		}
		if err = taskProvider.Input(providers.NewDefaultInput(pipeline.RootID,
			curTask.ID,
			curTask.Input,
			curTask.Env,
		)); err != nil {
			err = errors.Wrapf(err, "taskProvider %s check input err", task.Action)
			return
		}
		// when run task,refresh the task status value before making judgment
		taskState := curTask.RefreshState()
		if !taskState.IsReady() {
			w.rollback(taskProvider, curTask.Status)
			if taskState.IsSkip() {
				pipeline.AddSucceeded()
				continue
			}
			return
		}
		// 执行task主体
		if err = w.runWithRetry(curTask, taskProvider.Run, options.GetRetryGapSecond(curTask.Env)); err != nil {
			err = errors.Wrapf(err, "taskProvider %s run err", task.Action)
			return
		}
		// 处理task输出结果
		var output, context map[string]interface{}
		context, output, err = taskProvider.Output()
		if err != nil {
			err = errors.Wrapf(err, "taskProvider %s output err", task.Action)
			return
		}
		pipeCtx = helper.UnsafeMergeMap(pipeCtx, context)
		curTask.Context = pipeCtx
		curTask.Output = output
		if err = task.Succeed(); err != nil {
			err = errors.Wrapf(err, "taskProvider %s Succeed err", task.Action)
			return
		}
		pipeline.AddSucceeded()
		w.rollback(taskProvider, curTask.Status)
	}
}

// runWithRetry call the function and try again
func (w *DefaultWorkerPool) runWithRetry(task *dao.PipelineTask, fn func(int) error, sleep time.Duration) (err error) {
	for i := 1; i <= task.Retry; i++ {
		if err = fn(int(task.Retries)); err != nil {
			err = errors.Wrapf(err, "run err with retry nu:%d", i)
			time.Sleep(time.Duration(i) * sleep)
		} else {
			if task.Retries > 0 {
				task.Retries++
			}
			return
		}
		task.Retries++
	}
	return
}

func (w *DefaultWorkerPool) rollback(provider providers.TaskProvider, status string) {
	if provider == nil {
		return
	}
	rollback := providers.GetRollback(provider)
	if rollback != nil {
		rollback.Rollback(status)
	}
	providers.ReSet(provider)
}

func (w *DefaultWorkerPool) Quit() {
	w.stopOnce.Do(func() {
		w.isQuit = true
		close(w.pipePool)
	})
}

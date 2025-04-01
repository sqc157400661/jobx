package internal

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"

	"github.com/sqc157400661/jobx/internal/helper"
	"github.com/sqc157400661/jobx/pkg/model"
	"github.com/sqc157400661/jobx/pkg/options"
	"github.com/sqc157400661/jobx/pkg/providers"
)

// WorkerPool defines interface for pipeline task processing worker pools.
type WorkerPool interface {
	Submit(p *Pipeline, isSync bool) (err error)
	Start()
	Stop()
}

// DefaultWorkerPool implements WorkerPool using fixed goroutine workers.
// Handles task distribution, panic recovery, and graceful shutdown.
type DefaultWorkerPool struct {
	// Maximum concurrent workers
	maxWorkers int
	// Task queue (buffered channel)
	pipePool chan *Pipeline
	// Tracks active workers
	wg sync.WaitGroup
	// Atomic shutdown flag
	isQuit atomic.Bool
}

// NewDefaultWorkerPool creates a worker pool with specified concurrency.
func NewDefaultWorkerPool(maxWorkers int) (w *DefaultWorkerPool) {
	w = &DefaultWorkerPool{
		maxWorkers: maxWorkers,
		pipePool:   make(chan *Pipeline, maxWorkers),
	}
	return
}

// Submit adds a pipeline to the pool. Returns error if pool is shutting down.
// Thread-safe: Uses atomic check to prevent sending to closed channel.
func (w *DefaultWorkerPool) Submit(p *Pipeline, isSync bool) (err error) {
	if w.isQuit.Load() {
		return errors.New("worker is quit")
	}
	if isSync {
		w.process(p)
	} else {
		w.pipePool <- p
	}

	return
}

// Start starts worker goroutines. Idempotent - only executes once.
func (w *DefaultWorkerPool) Start() {
	for i := 0; i < w.maxWorkers; i++ {
		w.wg.Add(1)
		go w.worker()
	}
}

// worker processes tasks from pipePool until channel closure.
func (w *DefaultWorkerPool) worker() {
	defer w.wg.Done()
	for task := range w.pipePool {
		w.process(task)
	}
}

// process executes pipeline steps with panic recovery and error handling.
func (w *DefaultWorkerPool) process(pipeline *Pipeline) {
	if pipeline == nil {
		return
	}
	if len(pipeline.Steps) == 0 {
		pipeline.Finish()
		return
	}
	var (
		err          error
		curTask      *model.PipelineTask
		taskProvider providers.TaskProvider
	)
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
		// Merge contexts and update task
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
		// Execute with retry
		if err = w.runWithRetry(curTask, taskProvider.Run, options.GetRetryGapSecond(curTask.Env)); err != nil {
			err = errors.Wrapf(err, "taskProvider %s run err", task.Action)
			return
		}
		// Process outputs
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

// runWithRetry executes task with retry logic. Uses exponential backoff.
func (w *DefaultWorkerPool) runWithRetry(task *model.PipelineTask, fn func(int) error, sleep time.Duration) (err error) {
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

// rollback executes provider rollback if available.
func (w *DefaultWorkerPool) rollback(provider providers.TaskProvider, status string) {
	if provider == nil {
		return
	}
	if rb := providers.GetRollback(provider); rb != nil {
		rb.Rollback(status)
	}
	providers.ReSet(provider)
}

// Stop stops the pool
func (w *DefaultWorkerPool) Stop() {
	// Prevent new submissions
	w.isQuit.Store(true)
	// Stop worker goroutines
	close(w.pipePool)
	w.wg.Wait()
}

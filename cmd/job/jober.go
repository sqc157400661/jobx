package job

import (
	"context"
	"fmt"
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/internal"
	"github.com/sqc157400661/jobx/pkg/dao"
	"github.com/sqc157400661/jobx/pkg/errors"
	"github.com/sqc157400661/jobx/pkg/options"
	"github.com/sqc157400661/jobx/pkg/providers"
	"time"

	"strings"
)

// job拥有者
type Jober struct {
	Job         *dao.Job
	Level       int
	Pipeline    *internal.Pipeline
	root        *Jober
	child       *Jober
	Tokens      []string // 令牌
	err         error
	sync        bool
	syncTimeOut int
}

func NewJober(name, owner, tenant string, opts ...options.JobOptionFunc) *Jober {
	o := options.DefaultJobOptions
	for _, opt := range opts {
		opt(&o)
	}
	job := &Jober{
		Job: &dao.Job{
			Name:   name,
			Desc:   o.Desc,
			Pause:  o.Pause,
			Input:  o.Input,
			Env:    o.Env,
			BizID:  o.BizId,
			Locker: o.PreLockUid,
			Owner:  owner,
			Tenant: tenant,
			State: dao.State{
				Phase:  config.PhaseReady,
				Status: config.StatusPending,
			},
		},
		sync:        o.Sync,
		syncTimeOut: o.SyncTimeOut,
		Tokens:      o.Tokens,
		Level:       0,
	}
	return job
}

func (j *Jober) AddPipeline(name string, action string, opts ...options.JobOptionFunc) *Jober {
	o := options.DefaultJobOptions
	for _, opt := range opts {
		opt(&o)
	}
	if j.Err() != nil {
		return j
	}
	// 判断是否在全局Provider中
	if !providers.Has(action) {
		j.err = errors.NoProvider()
		return j
	}
	if j.Pipeline == nil {
		j.Pipeline = &internal.Pipeline{}
	}
	t := &dao.PipelineTask{
		Name:   name,
		Action: action,
		Desc:   providers.GetDesc(action, o.Desc),
		Pause:  o.Pause,
		Input:  internal.UnsafeMergeMap(o.Input, j.Job.Input),
		Env:    internal.UnsafeMergeMap(o.Env, j.Job.Env),
		Retry:  o.Retry,
		State: dao.State{
			Phase:  config.PhaseReady,
			Status: config.StatusPending,
		},
	}
	j.Pipeline.Tasks = append(j.Pipeline.Tasks, t)
	j.Job.Runnable = config.RunnableYes
	return j
}

func (j *Jober) AddJob(job *Jober) *Jober {
	if j.Pipeline != nil {
		job.err = errors.New(errors.ParamError, "jober can not add,because higher level jober already has pipeline")
	}
	if job.Err() != nil {
		return job
	}
	if j.root == nil {
		job.root = j
	} else {
		job.root = j.root
	}
	job.Job.BizID = "" // 子任务无效参数
	job.Tokens = nil   // 子任务无效参数
	job.Level = j.Level + 1
	job.Job.Input = internal.UnsafeMergeMap(job.Job.Input, j.Job.Input)
	job.Job.Env = internal.UnsafeMergeMap(job.Job.Env, j.Job.Env)
	if job.Level > config.MaxJobLevel {
		job.err = errors.New(errors.ParamError, "job level exceeds max limit")
	}
	j.child = job
	return job
}

func (j *Jober) Exec() (err error) {
	sess := dao.JFDb.NewSession()
	defer sess.Close()
	defer func() {
		if err != nil {
			_ = sess.Rollback()
		} else {
			_ = sess.Commit()
		}
	}()
	if j.Err() != nil {
		return j.Err()
	}
	var root *Jober
	root = j.root
	if j.root == nil {
		root = j
	}
	if root == nil || root.Job == nil {
		return errors.New(errors.ParamError, "root not available")
	}
	if root.Job.BizID != "" {
		var has bool
		var existJob dao.Job
		existJob, has, err = dao.GetJobByBizId(root.Job.BizID)
		if err != nil {
			return err
		}
		if has {
			j.Job = &existJob
			return errors.BIZConflict(fmt.Sprintf("biz:%s rootID:%d", root.Job.BizID, existJob.RootID))
		}
	}
	if len(root.Tokens) > 0 {
		var rootId int
		rootId, err = dao.CheckTokens(root.Tokens)
		if err != nil {
			j.Job.ID = rootId
			return
		}
	}
	_ = sess.Begin()
	var pre *Jober
	for jober := root; jober != nil; jober = jober.child {
		// save job
		if pre != nil {
			jober.Job.ParentID = pre.Job.ID
			jober.Job.Path = strings.TrimLeft(pre.Job.Path+fmt.Sprintf(",%d", pre.Job.ID), ",")
		}
		pre = jober
		jober.Job.RootID = root.Job.ID
		jober.Job.Level = jober.Level
		jober.Job.Phase = config.PhaseInit
		if jober.Job.ID == 0 {
			err = jober.Job.Save()
			if err != nil {
				return
			}
		}
		// save Tasks
		if jober.Pipeline == nil {
			continue
		}
		for _, task := range jober.Pipeline.Tasks {
			task.JobID = jober.Job.ID
			if task.Env == nil {
				task.Env = map[string]interface{}{"rootID": root.Job.ID}
			} else {
				task.Env["rootID"] = root.Job.ID
			}
			err = task.Save()
			if err != nil {
				return
			}
		}
		jober.Job.Phase = config.PhaseReady
		err = jober.Job.Update()
		if err != nil {
			return
		}
	}
	err = dao.CreateTokens(root.Job.ID, root.Tokens)
	if err != nil {
		return
	}
	if root.sync {
		return WaitJob(root.Job.ID, root.syncTimeOut)
	}
	return nil
}

func (j *Jober) Err() error {
	return j.err
}

func WaitJob(jobId int, timeout int) error {
	if timeout == 0 {
		timeout = 5
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel() // 确保释放资源
	// 等待所有任务完成或超时
	done := make(chan struct{})
	go waitJob(ctx, jobId, done)
	select {
	case <-ctx.Done(): // 超时
		return errors.WaitTimeout()
	case <-done: // 所有任务完成
		return nil
	}
}

func waitJob(ctx context.Context, jobId int, done chan struct{}) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered waitJob:", r)
			}
		}()
	}()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop() // 确保在函数返回时停止定时器
	var job dao.Job
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, _ = dao.JFDb.ID(jobId).Get(&job)
			if job.State.IsSuccess() {
				close(done)
				return
			}
		}
	}
}

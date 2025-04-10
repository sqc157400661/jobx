package job

import (
	"fmt"
	"github.com/sqc157400661/jobx/pkg/mysql"
	"time"

	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/pkg/errors"
	"github.com/sqc157400661/jobx/pkg/model"
)

type JobDefinition struct {
	JobDef JobDef `yaml:"jobDef"`
}

type JobDef struct {
	Name      string                 `yaml:"name"`
	Desc      string                 `yaml:"desc"`
	Owner     string                 `yaml:"owner"`
	AppName   string                 `yaml:"appName"`
	Tenant    string                 `yaml:"tenant"`
	BizID     string                 `yaml:"bizId"`
	Pipelines []Pipeline             `yaml:"pipelines,omitempty"`
	Locker    string                 `yaml:"locker,omitempty"`
	Jobs      []JobDef               `yaml:"jobs,omitempty"`
	Input     map[string]interface{} `yaml:"input,omitempty"`
	Env       map[string]interface{} `yaml:"env,omitempty"`
	Pause     bool                   `yaml:"pause,omitempty"`
}

type Pipeline struct {
	Name      string                 `yaml:"name"`
	Desc      string                 `yaml:"desc"`
	Action    string                 `yaml:"action"`
	Pause     bool                   `yaml:"pause,omitempty"`
	Rollback  bool                   `yaml:"rollback,omitempty"`
	RetryNum  int                    `yaml:"retryNum,omitempty"`
	Condition []string               `yaml:"condition,omitempty"`
	Input     map[string]interface{} `yaml:"input,omitempty"`
	Env       map[string]interface{} `yaml:"env,omitempty"`
}

func SaveJobsFromDef(parentDef JobDef) (rootId int, err error) {
	sess := mysql.DB().NewSession()
	err = sess.Begin()
	defer func() {
		if err != nil {
			_ = sess.Rollback()
		} else {
			_ = sess.Commit()
		}
	}()
	if err != nil {
		return
	}
	type queueItem struct {
		def    *JobDef
		parent *model.Job
	}
	parentJob, err := createJobFromDef(&parentDef, &model.Job{})
	if err != nil {
		return
	}
	rootId = parentJob.ID
	parentJob.RootID = parentJob.ID
	parentJob.Level = 1
	err = parentJob.Update()
	if err != nil {
		return
	}
	var queue []queueItem
	// 初始化队列，加入第一层子任务
	for i := range parentDef.Jobs {
		queue = append(queue, queueItem{
			def:    &parentDef.Jobs[i],
			parent: parentJob,
		})
	}

	var childJob *model.Job
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// 处理当前任务
		childJob, err = createJobFromDef(current.def, current.parent)
		if err != nil {
			return
		}

		// 更新层级信息
		childJob.Level = current.parent.Level + 1
		childJob.Path = fmt.Sprintf("%s,%d", current.parent.Path, childJob.ID)

		if err = childJob.Save(); err != nil {
			// todo define error
			return 0, errors.NewDBError(fmt.Sprintf("failed to update child job %s: %w", current.def.Name, err))
		}

		// 将子任务加入队列
		for i := range current.def.Jobs {
			queue = append(queue, queueItem{
				def:    &current.def.Jobs[i],
				parent: childJob,
			})
		}
	}
	return
}

func createJobFromDef(def *JobDef, parent *model.Job) (*model.Job, error) {
	job := &model.Job{
		AppName:  def.AppName,
		Tenant:   def.Tenant,
		Input:    def.Input,
		Owner:    def.Owner,
		Env:      def.Env,
		Name:     def.Name,
		Desc:     def.Desc,
		Level:    parent.Level + 1,
		Pause:    boolToInt8(def.Pause),
		RootID:   parent.RootID,
		BizID:    def.BizID,
		ParentID: parent.ID,
		CreateAt: time.Now().Unix(),
		UpdateAt: time.Now().Unix(),
		State: model.State{
			Phase:  config.PhaseReady,
			Status: config.StatusPending,
		},
	}
	if parent.ID > 0 {
		job.Path = fmt.Sprintf("%s,%d", parent.Path, parent.ID)
	}
	var err error
	// Determine if this job is runnable (has pipelines directly)
	job.Runnable = boolToInt8(len(def.Pipelines) > 0 && len(def.Jobs) == 0)

	// Save to database
	if err = job.Save(); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// If this job is runnable, create pipeline tasks
	if job.Runnable == 1 {
		if err = createPipelineTasks(def.Pipelines, job.ID); err != nil {
			return nil, fmt.Errorf("failed to create pipeline tasks: %w", err)
		}
	}

	return job, nil
}

func createPipelineTasks(pipelines []Pipeline, jobID int) error {
	for _, pipeline := range pipelines {
		task := &model.PipelineTask{
			JobID:    jobID,
			Name:     pipeline.Name,
			Action:   pipeline.Action,
			Desc:     pipeline.Desc,
			Pause:    boolToInt8(pipeline.Pause),
			Retry:    pipeline.RetryNum,
			Input:    pipeline.Input,
			Env:      pipeline.Env,
			CreateAt: time.Now().Unix(),
			UpdateAt: time.Now().Unix(),
			State: model.State{
				Phase:  config.PhaseReady,
				Status: config.StatusPending,
			},
		}

		if err := task.Save(); err != nil {
			return fmt.Errorf("failed to create job task: %w", err)
		}
	}
	return nil
}

// Helper functions

func boolToInt8(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

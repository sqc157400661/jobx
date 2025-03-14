package biz

import (
	"github.com/sqc157400661/jobx/api/types"
	"github.com/sqc157400661/jobx/cmd/job"
	"github.com/sqc157400661/jobx/pkg/options"
)

// Create a job flow
func Create(req types.CreateJober) (err error) {
	var j *job.Jober
	j, err = create(&req)
	if err != nil {
		return
	}
	return j.Exec()
}

func create(c *types.CreateJober) (j *job.Jober, err error) {
	j = job.NewJober(
		c.Name,
		c.Owner,
		c.Tenant,
		options.BizId(c.BizId),
		options.JobDesc(c.Desc),
		options.Pause(c.Pause),
		options.JobInput(c.Input),
		options.JobEnv(c.Env),
	)
	if len(c.Pipelines) > 0 {
		err = createPipelines(j, c.Pipelines)
		return
	}
	for _, cj := range c.ChildJobs {
		var childJober *job.Jober
		childJober, err = create(cj)
		if err != nil {
			return
		}
		j.AddJob(childJober)
	}
	return
}

func createPipelines(j *job.Jober, pipelines []*types.CreatePipeline) (err error) {
	for _, pip := range pipelines {
		j.AddPipeline(
			pip.Name,
			pip.Action,
			options.JobDesc(pip.Desc),
			options.RetryNum(pip.Retry),
			options.Pause(pip.Pause),
			options.JobInput(pip.Input),
			options.JobEnv(pip.Env),
		)
		if j.Err() != nil {
			return j.Err()
		}
	}
	return
}

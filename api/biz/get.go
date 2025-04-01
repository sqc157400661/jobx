package biz

import (
	"github.com/sqc157400661/jobx/api/types"
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/pkg/model"
)

func Get(id int) (res types.JobResult, err error) {
	var job model.Job
	var hasJob bool
	job, hasJob, err = model.GetJobById(id)
	if !hasJob || err != nil {
		return
	}
	res.Job = &job
	var childJobs []*model.Job
	childJobs, err = model.GetChildJobsById(id)
	if err != nil {
		return
	}
	res.ChildJobs = childJobs
	if job.Runnable == config.RunnableYes {
		var tasks []*model.PipelineTask
		tasks, err = model.GetPipelineTasksByJobId(id)
		if err != nil {
			return
		}
		res.Task = tasks
	}
	return
}

func GetByBizID(bid string) (job model.Job, has bool, err error) {
	return model.GetJobByBizId(bid)
}

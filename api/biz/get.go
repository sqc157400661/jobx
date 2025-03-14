package biz

import (
	"github.com/sqc157400661/jobx/api/types"
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/pkg/dao"
)

func Get(id int) (res types.JobResult, err error) {
	var job dao.Job
	var hasJob bool
	job, hasJob, err = dao.GetJobById(id)
	if !hasJob || err != nil {
		return
	}
	res.Job = &job
	var childJobs []*dao.Job
	childJobs, err = dao.GetChildJobsById(id)
	if err != nil {
		return
	}
	res.ChildJobs = childJobs
	if job.Runnable == config.RunnableYes {
		var tasks []*dao.PipelineTask
		tasks, err = dao.GetPipelineTasksByJobId(id)
		if err != nil {
			return
		}
		res.Task = tasks
	}
	return
}

func GetByBizID(bid string) (job dao.Job, has bool, err error) {
	return dao.GetJobByBizId(bid)
}

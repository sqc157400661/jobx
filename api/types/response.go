package types

import "github.com/sqc157400661/jobx/pkg/dao"

type JobResult struct {
	Job       *dao.Job
	ChildJobs []*dao.Job
	Task      []*dao.PipelineTask
}

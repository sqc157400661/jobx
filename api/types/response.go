package types

import "github.com/sqc157400661/jobx/pkg/model"

type JobResult struct {
	Job       *model.Job
	ChildJobs []*model.Job
	Task      []*model.PipelineTask
}

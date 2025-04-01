package storage

import (
	"github.com/go-xorm/xorm"
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/pkg/model"
)

type JobStorage struct {
	Orm *xorm.Session
}

func NewJobStorage(sess *xorm.Session) *JobStorage {
	return &JobStorage{Orm: sess}
}

func (s *JobStorage) SaveJob(job *model.Job) (err error) {
	_, err = s.Orm.InsertOne(job)
	return
}

func (s *JobStorage) UpdateJob(job *model.Job) (err error) {
	_, err = s.Orm.Update(job, &model.Job{ID: job.ID})
	return
}

func (s *JobStorage) MarkJobRunning(job *model.Job) (err error) {
	job.State.Phase = config.PhaseRunning
	_, err = s.Orm.Update(job, &model.Job{ID: job.ID, State: model.State{Phase: config.PhaseReady}})
	return
}

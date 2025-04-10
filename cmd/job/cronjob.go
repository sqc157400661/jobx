package job

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/pkg/errors"
	"github.com/sqc157400661/jobx/pkg/model"
	"github.com/sqc157400661/jobx/pkg/options/cronopt"
)

func NewCronjob(spec, name, owner string, opts ...cronopt.CronOptionFunc) (*Cronjob, error) {
	o := cronopt.DefaultOption
	for _, opt := range opts {
		opt(&o)
	}
	cronParseOption := cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow
	if o.SecondEnable {
		cronParseOption = cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow
	}
	parser := cron.NewParser(cronParseOption)
	_, err := parser.Parse(spec)
	if err != nil {
		return nil, errors.NewParamError("invalid spec").Wrap(err)
	}
	return &Cronjob{
		spec:  spec,
		name:  name,
		owner: owner,
		opt:   o,
	}, nil
}

func (c *Cronjob) ExecJob(job *Jober) (err error) {
	if job == nil {
		return nil
	}
	sess := mysql.DB().NewSession()
	//jobStorage := storage.NewJobStorage(sess)
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return errors.NewDBError("begin err").Wrap(err)
	}
	defer func() {
		if err != nil {
			_ = sess.Rollback()
		} else {
			_ = sess.Commit()
		}
	}()
	job.Name = fmt.Sprintf("%s_%s", c.name, job.Name)
	job.AppName = c.opt.AppName
	job.Tenant = c.opt.Tenant
	job.Owner = c.owner
	existCronJob, err := model.IsCronJobExist(c.opt.Tenant, c.opt.AppName, c.name)
	if err != nil {
		return
	}
	if existCronJob.ID > 0 {
		return errors.ErrCronJobExist
	}
	jobDef, err := model.GetJobDefByName(job.Name)
	if err != nil {
		return err
	}
	if jobDef.ID > 0 {
		return errors.ErrJobDefExist
	} else {
		var yamlConf string
		yamlConf, err = job.ToYAML()
		if err != nil {
			return err
		}
		jobDef = &model.JobDefinition{
			Name:     job.Name,
			AppName:  c.opt.AppName,
			Tenant:   c.opt.Tenant,
			YamlConf: yamlConf,
		}
		if err = jobDef.Save(); err != nil {
			return err
		}
	}
	cronJob := model.JobCron{
		Name:           c.name,
		Owner:          c.owner,
		Status:         config.CronStatusValid,
		Spec:           c.spec,
		ExecType:       config.JobExecType,
		ExecContent:    job.Name,
		Tenant:         c.opt.Tenant,
		AppName:        c.opt.AppName,
		CurrencyPolicy: c.opt.CurrencyPolicy,
	}
	return cronJob.Save()
}

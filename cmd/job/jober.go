package job

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"

	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/internal/helper"
	"github.com/sqc157400661/jobx/internal/job"
	"github.com/sqc157400661/jobx/pkg/errors"
	"github.com/sqc157400661/jobx/pkg/model"
	"github.com/sqc157400661/jobx/pkg/options/cronopt"
	"github.com/sqc157400661/jobx/pkg/options/jobopt"
	"github.com/sqc157400661/jobx/pkg/providers"
)

type Cronjob struct {
	name  string
	spec  string
	owner string
	opt   cronopt.CronOptions
}

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
	sess := model.DB().NewSession()
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

func YAMLToJob(yamlContent []byte) (int, error) {
	var jobDef job.JobDefinition
	var err error
	// Unmarshal YAML
	if err = yaml.Unmarshal(yamlContent, &jobDef); err != nil {
		return 0, errors.NewParamError("failed to unmarshal YAML").Wrap(err)
	}
	var rootId int
	// Process child jobs recursively
	if rootId, err = job.SaveJobsFromDef(jobDef.JobDef); err != nil {
		return 0, errors.NewParamError("failed to process child jobs").Wrap(err)
	}
	return rootId, nil
}

// Jober 构建Job定义的链式调用结构体
type Jober struct {
	job.JobDef
	level       int      `yaml:"-"`
	tokens      []string `yaml:"-"`
	err         error    `yaml:"-"`
	sync        bool     `yaml:"-"`
	syncTimeOut int      `yaml:"-"`
	JobId       int      `yaml:"-"`
	engineType  string   `yaml:"-"`
}

// NewJob 创建新的Job实例
func NewJob(name, owner string, opts ...jobopt.JobOptionFunc) *Jober {
	o := jobopt.DefaultJobOptions
	for _, opt := range opts {
		opt(&o)
	}
	return &Jober{
		JobDef: job.JobDef{
			Name:  name,
			Desc:  o.Desc,
			Env:   o.Env,
			Input: o.Input,
			//Pause: o.Pause,
			Locker:  o.PreLockUid,
			BizID:   o.BizId,
			Owner:   owner,
			AppName: o.AppName,
			Tenant:  o.Tenant,
		},
		level:       1,
		sync:        o.Sync,
		syncTimeOut: o.SyncTimeOut,
		tokens:      o.Tokens,
	}
}

func (j *Jober) AddJob(job *Jober) *Jober {
	if j.Pipelines != nil {
		job.err = errors.NewParamError("jober can not add,because higher level jober already has pipeline")
	}
	if job.err != nil {
		return job
	}
	job.level = j.level + 1
	job.BizID = ""   // 子任务无效参数
	job.tokens = nil // 子任务无效参数
	job.Input = helper.UnsafeMergeMap(job.Input, j.Input)
	job.Env = helper.UnsafeMergeMap(job.Env, j.Env)
	if job.level > config.MaxJobLevel {
		job.err = errors.NewParamError("job level exceeds max limit")
	}
	j.Jobs = append(j.Jobs, job.JobDef)
	return job
}

func (j *Jober) AddPipeline(name string, action string, opts ...jobopt.JobOptionFunc) *Jober {
	o := jobopt.DefaultJobOptions
	for _, opt := range opts {
		opt(&o)
	}
	if j.err != nil {
		return j
	}
	// 判断是否在全局Provider中
	if !providers.Has(action) {
		j.err = errors.ErrNoProvider
		return j
	}
	t := job.Pipeline{
		Name:   name,
		Action: action,
		Desc:   providers.GetDesc(action, o.Desc),
		//Pause:  o.Pause,
		Input:    helper.UnsafeMergeMap(o.Input, j.Input),
		Env:      helper.UnsafeMergeMap(o.Env, j.Env),
		RetryNum: o.Retry,
	}
	j.Pipelines = append(j.Pipelines, t)
	return j
}

func (j *Jober) Exec() (err error) {
	if j == nil {
		return
	}
	if j.err != nil {
		return j.err
	}
	sess := model.DB().NewSession()
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
	if err = j.checkBizID(); err != nil {
		return err
	}
	if err = j.checkTokens(); err != nil {
		return err
	}
	var rootId int
	if rootId, err = job.SaveJobsFromDef(j.JobDef); err != nil {
		return
	}
	err = model.CreateTokens(rootId, j.tokens)
	if err != nil {
		return
	}
	if j.sync {
		return WaitJob(rootId, j.syncTimeOut)
	}
	return nil
}

func (j *Jober) checkBizID() (err error) {
	if j.BizID == "" {
		return nil
	}
	var has bool
	var existJob model.Job
	existJob, has, err = model.GetJobByBizId(j.BizID)
	if err != nil {
		return err
	}
	if has {
		j.JobId = existJob.ID
		return errors.BIZConflict(fmt.Sprintf("biz:%s rootID:%d", j.BizID, existJob.RootID))
	}
	return nil
}

func (j *Jober) checkTokens() (err error) {
	if len(j.tokens) == 0 {
		return nil
	}
	var rootId int
	rootId, err = model.CheckTokens(j.tokens)
	if err != nil {
		j.JobId = rootId
		return
	}
	return nil
}

// ToYAML 生成YAML字符串
func (j *Jober) ToYAML() (string, error) {
	def := &job.JobDefinition{
		JobDef: j.JobDef,
	}
	yamlData, err := yaml.Marshal(def)
	if err != nil {
		return "", fmt.Errorf("failed to marshal job definition: %w", err)
	}
	return string(yamlData), nil
}

// AddInput 添加入参
func (j *Jober) AddInput(key string, value interface{}) *Jober {
	if j.Input == nil {
		j.Input = make(map[string]interface{})
	}
	j.Input[key] = value
	return j
}

// AddEnv 添加环境变量
func (j *Jober) AddEnv(key string, value interface{}) *Jober {
	if j.Env == nil {
		j.Env = make(map[string]interface{})
	}
	j.Env[key] = value
	return j
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
		return errors.ErrWaitJobTimeout
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
	ticker := time.NewTicker(200 * time.Microsecond)
	defer ticker.Stop() // 确保在函数返回时停止定时器
	var job model.Job
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_, _ = model.DB().ID(jobId).Get(&job)
			if job.State.IsSuccess() {
				close(done)
				return
			}
		}
	}
}

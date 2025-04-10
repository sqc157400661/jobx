package collector

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/go-xorm/xorm"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"

	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/internal/names"
	"github.com/sqc157400661/jobx/pkg/model"
)

const (
	checkUndoJobsSqlTmpl = "locker=? and phase !=? and parent_id=0 and tenant=? and app_name=?"
	getStealJobSqlTmpl   = "locker=? and phase =? and parent_id=0 and tenant=? and app_name=?"
	stealJobsSqlTmpl     = `update job set locker=? where  parent_id=0 and (locker='' or locker=?) and phase =? and tenant=? and app_name=? order by id asc limit ?`
	releaseJobSqlTmpl    = `update job set locker='',phase =?  where locker=? and id =?`
	stealCronJobsSqlTmpl = `update job_cron set locker=? where locker='' and status =? and tenant=? and app_name=? order by id asc limit ?`
)

type Collector interface {
	StealJobs() (jobs []*model.Job, err error)
	ReleaseJobs() (err error)
	ReleaseJobByID(jobId int) (err error)
	StealCronJobs() (jobs []*model.JobCron, err error)
}

type DefaultCollector struct {
	stealLen  int
	serverUid string
	tenant    string
	appName   string
	engine    *xorm.Engine
}

func NewDefaultCollector(engine *xorm.Engine, serverUid, tenant, appName string, len int) (collector *DefaultCollector) {
	if len <= 0 {
		len = 1
	}
	return &DefaultCollector{
		tenant:    tenant,
		appName:   appName,
		serverUid: serverUid,
		stealLen:  len,
		engine:    engine,
	}
}

// 启动初始化的时候检查未完成的job队列
func (c *DefaultCollector) loadCheckUndoJobs() (jobs []model.Job, err error) {
	err = c.engine.Where(checkUndoJobsSqlTmpl, c.serverUid, config.PhaseTerminated, c.tenant, c.appName).Find(&jobs)
	return
}

func (c *DefaultCollector) StealJobs() (jobs []*model.Job, err error) {
	var num int64
	num, err = c.steal()
	if err != nil || num == 0 {
		return
	}
	err = c.engine.Where(getStealJobSqlTmpl, c.serverUid, config.PhaseReady, c.tenant, c.appName).Find(&jobs)
	return
}

func (c *DefaultCollector) ReleaseJobs() (err error) {
	var jobs []model.Job
	err = c.engine.In("phase", []string{config.PhaseReady, config.PhaseRunning}).Where("locker=?", c.serverUid).Find(&jobs)
	if err != nil {
		return
	}
	// 依次更新状态并解除锁定
	for _, v := range jobs {
		err = c.ReleaseJobByID(v.ID)
		if err != nil {
			err = errors.Wrapf(err, "uid:%s unlock err id:%d", c.serverUid, v.ID)
			return
		}
	}
	return
}

// AddJobLocker 对job添加locker
func (c *DefaultCollector) steal() (lockedNum int64, err error) {
	var res sql.Result
	// CAS
	res, err = c.engine.Exec(stealJobsSqlTmpl, c.serverUid, names.PreLockKey(c.serverUid), config.PhaseReady, c.tenant, c.appName, c.stealLen)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (c *DefaultCollector) ReleaseJobByID(jobId int) (err error) {
	_, err = c.engine.Exec(releaseJobSqlTmpl, config.PhaseReady, c.serverUid, jobId)
	return
}

func (c *DefaultCollector) StealCronJobs() (jobs []*model.JobCron, err error) {
	var num int64
	num, err = c.stealCronJob()
	if err != nil || num == 0 {
		return
	}
	err = c.engine.Where("locker=?", c.serverUid).Find(&jobs)
	return
}

func (c *DefaultCollector) stealCronJob() (lockedNum int64, err error) {
	var res sql.Result
	// CAS
	res, err = c.engine.Exec(stealCronJobsSqlTmpl, c.serverUid, config.CronStatusValid, c.tenant, c.appName, 2)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

/* ----------------------------------------------------------------------*/
type CronTrigger interface {
	Load(jobs []*model.JobCron)
	Start()
	Stop()
}

type DefaultCronTrigger struct {
	serverUid string
	cron      *cron.Cron
	mutex     sync.RWMutex
	cronJobs  map[int64]*model.JobCron
	entryIDs  sets.Int
}

func NewDefaultCronTrigger() *DefaultCronTrigger {
	cronTrigger := &DefaultCronTrigger{
		cronJobs: make(map[int64]*model.JobCron),
		entryIDs: sets.NewInt(),
	}
	return cronTrigger
}

func (r *DefaultCronTrigger) Load(jobs []*model.JobCron) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, job := range jobs {
		r.cronJobs[job.ID] = job
	}
}

func (r *DefaultCronTrigger) Start() {
	r.cron = cron.New(cron.WithSeconds())
	// 添加同步任务并处理潜在错误
	if _, err := r.cron.AddFunc("*/1 * * * * *", r.syncCronJobs); err != nil {
		klog.Errorf("Failed to add sync job: %v", err)
	}
	r.cron.Start()
}

func (r *DefaultCronTrigger) syncCronJobs() {
	for _, c := range r.cronJobs {
		if err := r.processCronJob(c); err != nil {
			klog.Errorf("Process cron job failed: %v cronID: %d", err, c.ID)
		}
	}
}

func (r *DefaultCronTrigger) processCronJob(c *model.JobCron) error {
	switch {
	case c.IsDeleting():
		return r.handleDeletingJob(c)
	case c.IsRebooting():
		return r.handleRebootingJob(c)
	case c.IsRunning():
		return r.healthCheckJob(c)
	case c.NeedsScheduling():
		return r.scheduleJob(c)
	}
	return nil
}

func (r *DefaultCronTrigger) handleDeletingJob(c *model.JobCron) error {
	// 优先删除数据库记录
	if _, err := mysql.DB().ID(c.ID).Delete(c); err != nil {
		return fmt.Errorf("delete cronjob failed: %w", err)
	}
	r.remove(c.EntryID)
	return nil
}

func (r *DefaultCronTrigger) handleRebootingJob(c *model.JobCron) error {
	if err := c.MarkValid(); err != nil {
		return fmt.Errorf("mark valid failed: %w", err)
	}
	r.remove(c.EntryID)
	return nil
}

func (r *DefaultCronTrigger) healthCheckJob(c *model.JobCron) error {
	if r.entryIDs.Has(c.EntryID) {
		return c.HealthCheckOk()
	} else {
		return c.MarkRebooting()
	}
}

func (r *DefaultCronTrigger) scheduleJob(c *model.JobCron) error {
	if r.entryIDs.Has(c.EntryID) {
		r.remove(c.EntryID)
	}
	entryID, err := r.cron.AddJob(c.Spec, c)
	if err != nil {
		return fmt.Errorf("add job failed: %w", err)
	}

	if err = c.TouchEntryID(int(entryID)); err != nil {
		return fmt.Errorf("update entry ID failed: %w", err)
	}
	r.mutex.Lock()
	r.entryIDs.Insert(int(entryID))
	r.mutex.Unlock()
	if err = c.MarkRunning(); err != nil {
		return fmt.Errorf("mark running failed: %w", err)
	}
	return nil
}

func (r *DefaultCronTrigger) remove(entryID int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.entryIDs.Has(entryID) {
		r.cron.Remove(cron.EntryID(entryID))
		r.entryIDs.Delete(entryID)
	}
}

func (r *DefaultCronTrigger) Add(job *model.JobCron) error {
	if job.ExecContent == "" {
		return fmt.Errorf("ExecContent is empty ")
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.cronJobs[job.ID] != nil {
		return nil
	}
	r.cronJobs[job.ID] = job
	return nil
}

func (r *DefaultCronTrigger) Stop() {
	if r.cron != nil {
		r.cron.Stop()
	}
	for _, job := range r.cronJobs {
		_ = job.MarkValid()
	}
}

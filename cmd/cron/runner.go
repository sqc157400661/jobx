package cron

import (
	"errors"
	"github.com/robfig/cron/v3"
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/pkg/dao"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
	"sync"
)

type JobEntry interface {
	cron.Job
	Key() string
}
type Runner struct {
	uid       string
	cron      *cron.Cron
	startOnce sync.Once
	mutex     sync.RWMutex
	cronJobs  map[string]cron.Job
	entryIDs  sets.Int
}

func NewCronRunner(uid string, jobs ...JobEntry) (runner *Runner, err error) {
	if len(uid) == 0 {
		err = errors.New("no uid")
		return
	}
	runner = &Runner{
		uid:      uid,
		cronJobs: make(map[string]cron.Job),
		entryIDs: sets.NewInt(),
	}
	if len(jobs) > 0 {
		for _, job := range jobs {
			runner.cronJobs[job.Key()] = job
		}
	}
	return
}

func (r *Runner) Start() (err error) {
	// 初始化定时任务对象
	// Cron v3默认支持精确到分钟的cron表达式
	// cron.WithSeconds()表示指定支持精确到秒的表达式
	r.startOnce.Do(func() {
		r.cron = cron.New(cron.WithSeconds())
		// 每隔1分钟中check db
		_, err = r.cron.AddFunc("0 */1 * * * *", func() {

			var crons []*dao.JobCron
			crons, err = dao.CollectCron(r.uid)
			if err != nil {
				klog.Errorf("collectCron err %s,uid:%s", err.Error(), r.uid)
				return
			}
			var entryID cron.EntryID
			// 从数据库中获取已计划的任务并运行
			for _, c := range crons {
				// 判断是否在删除中
				if c.IsDeleted() {
					r.remove(c.EntryID)
					// 从数据库删除该定时任务数据
					_, err = dao.JFDb.ID(c.ID).Delete(c)
					if err != nil {
						klog.Errorf("delete cronjob err %s,uid:%s cronId:%d", err.Error(), r.uid, c.ID)
						continue
					}
				}
				// 判断是否在重启中
				if c.IsRebooting() {
					r.remove(c.EntryID)
					err = c.MarkValid()
					if err != nil {
						klog.Errorf("markValid err %s,uid:%s cronId:%d", err.Error(), r.uid, c.ID)
						continue
					}
				}
				// 判断是否在更新中
				if c.IsUpdating() {
					// 删除原定时任务
					r.remove(c.EntryID)
				}
				// 判断是否新执行或者已经在执行中
				if c.IsValid() || c.IsRunning() || c.IsUpdating() {
					if !r.entryIDs.Has(c.EntryID) {
						job, has := r.getCronJob(c.ExecContent)
						if !has {
							klog.Errorf("job not found uid:%s ExecContent:%s cronId:%d", r.uid, c.ExecContent, c.ID)
							continue
						}
						entryID, err = r.cron.AddJob(c.Spec, job)
						if err != nil {
							klog.Errorf("addJob err:%s uid:%s cronId:%d", err.Error(), r.uid, c.ID)
							continue
						}
						r.entryIDs.Insert(int(entryID))
						// 更新任务id
						err = c.TouchEntryID(int(entryID))
						if err != nil {
							klog.Errorf("touchEntryID err:% uid:%s cronId:%d", err.Error(), r.uid, c.ID)
							continue
						}
					}
					if !c.IsRunning() {
						err = c.MarkRunning()
						if err != nil {
							klog.Errorf("markRunning err %s,uid:%s cronId:%d", err.Error(), r.uid, c.ID)
						}
					}
				}
			}
		})
		if err != nil {
			return
		}
		//r.cron.Entry()
		// 启动Cron
		r.cron.Start()
	})
	return
}

func (r *Runner) getCronJob(name string) (job cron.Job, has bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	job, has = r.cronJobs[name]
	return
}

func (r *Runner) remove(entryID int) {
	if r.entryIDs.Has(entryID) {
		// 移除定时任务
		r.cron.Remove(cron.EntryID(entryID))
		r.entryIDs.Delete(int(entryID))
	}
	return
}

func (r *Runner) AddJobEntry(job JobEntry) {
	r.mutex.Lock()
	r.cronJobs[job.Key()] = job
	r.mutex.Unlock()
}

func (r *Runner) Add(spec, name string, job cron.Job) (id int64, err error) {
	r.mutex.Lock()
	r.cronJobs[name] = job
	r.mutex.Unlock()
	// 判断是否存在
	var cronData *dao.JobCron
	cronData, err = dao.GetCronByExecContent(name)
	if err != nil {
		return
	}
	if cronData.ID > 0 {
		id = cronData.ID
		return
	}
	cronData = &dao.JobCron{
		ExecType:    config.CronExecJobType,
		ExecContent: name,
		Spec:        spec,
		Locker:      r.uid,
		Status:      config.CronStatusValid,
	}
	// 将定时任务添加到数据库
	_, err = cronData.Add()
	if err != nil {
		return
	}
	id = cronData.ID
	return
}

func (r *Runner) Stop() {
	r.cron.Stop()
}

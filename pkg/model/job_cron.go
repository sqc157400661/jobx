package model

import (
	"time"

	"github.com/sqc157400661/jobx/config"
)

// JobCron
type JobCron struct {
	ID             int64     `gorm:"primaryKey;column:id" json:"id" xorm:"id pk autoincr"`
	Name           string    `gorm:"column:name" json:"name" xorm:"name"`                                     // 任务名称
	Owner          string    `gorm:"column:owner" json:"owner" xorm:"owner"`                                  // 任务归属人
	CurrencyPolicy string    `gorm:"column:currency_policy" json:"currency_policy" xorm:"currency_policy"`    // 并发策略
	EntryID        int       `gorm:"column:entry_id" json:"entry_id" xorm:"entry_id"`                         // 定时任务id
	Spec           string    `gorm:"column:spec" json:"spec" xorm:"spec"`                                     // 定时表达式
	ExecType       string    `gorm:"column:exec_type" json:"exec_type" xorm:"exec_type"`                      // 执行任务类型，如job、func、shell
	ExecContent    string    `gorm:"column:exec_content" json:"exec_content" xorm:"exec_content"`             // 执行任务内容
	Status         string    `gorm:"column:status" json:"status" xorm:"status"`                               // 状态
	AppName        string    `gorm:"column:app_name" json:"app_name" xorm:"app_name"`                         // 应用名称
	Tenant         string    `gorm:"column:tenant" json:"tenant" xorm:"tenant"`                               // tenant
	Locker         string    `gorm:"column:locker" json:"locker" xorm:"locker"`                               // 锁拥有者
	LastHealthTime time.Time `gorm:"column:last_health_time" json:"last_health_time" xorm:"last_health_time"` // 健康检查的时间
	CreateAt       int       `gorm:"column:create_at" json:"create_at" xorm:"created"`                        // 创建时间
	UpdateAt       int       `gorm:"column:update_at" json:"update_at" xorm:"updated"`                        // 更新时间
}

func (j *JobCron) Run() {
	// todo
}

func (j *JobCron) TableName() string {
	return "job_cron"
}

func (j *JobCron) Save() (err error) {
	_, err = DB().InsertOne(j)
	return
}

func (j *JobCron) TouchEntryID(id int) (err error) {
	if j == nil {
		return
	}
	j.EntryID = id
	_, err = DB().ID(j.ID).Update(j)
	return
}

func (j *JobCron) MarkRunning() (err error) {
	if j == nil {
		return
	}
	j.Status = config.CronStatusRunning
	_, err = DB().ID(j.ID).Update(j)
	return
}

func (j *JobCron) MarkValid() (err error) {
	if j == nil {
		return
	}
	j.Locker = ""
	j.Status = config.CronStatusValid
	_, err = DB().ID(j.ID).Update(j)
	return
}

func (j *JobCron) MarkUpdating() (err error) {
	if j == nil {
		return
	}
	j.Status = config.CronStatusUpdating
	_, err = DB().ID(j.ID).Update(j)
	return
}

func (j *JobCron) HealthCheckOk() (err error) {
	if j == nil {
		return
	}
	j.LastHealthTime = time.Now()
	_, err = DB().ID(j.ID).Update(j)
	return
}

func (j *JobCron) MarkDeleted() (err error) {
	if j == nil {
		return
	}
	j.Status = config.CronStatusDeleted
	_, err = DB().ID(j.ID).Update(j)
	return
}

func (j *JobCron) MarkRebooting() (err error) {
	if j == nil {
		return
	}
	j.Status = config.CronStatusRebooting
	_, err = DB().ID(j.ID).Update(j)
	return
}

// Add 添加cron任务
func (j *JobCron) Add() (id int64, err error) {
	if j == nil {
		return
	}
	id, err = DB().Insert(j)
	return
}

// Update 更新
func (j *JobCron) Update() (err error) {
	if j == nil {
		return
	}
	_, err = DB().ID(j.ID).Update(j)
	return
}

// IsRebooting 是否在重启中
func (j *JobCron) IsRebooting() bool {
	if j == nil {
		return false
	}
	return j.Status == config.CronStatusRebooting
}

// IsValid 是否在等待执行中
func (j *JobCron) IsValid() bool {
	if j == nil {
		return false
	}
	return j.Status == config.CronStatusValid
}

// IsRunning 是否在运行中
func (j *JobCron) IsRunning() bool {
	if j == nil {
		return false
	}
	return j.Status == config.CronStatusRunning
}

// IsDeleted 是否已经删除
func (j *JobCron) IsDeleted() bool {
	if j == nil {
		return false
	}
	return j.Status == config.CronStatusDeleted
}

// IsDeleting 是否正在删除
func (j *JobCron) IsDeleting() bool {
	if j == nil {
		return false
	}
	return j.Status == config.CronStatusDeleting
}

// IsUpdating 是否在更新状态
func (j *JobCron) IsUpdating() bool {
	if j == nil {
		return false
	}
	return j.Status == config.CronStatusUpdating
}

// NeedsScheduling 是否校验和更新状态
func (j *JobCron) NeedsScheduling() bool {
	if j == nil {
		return false
	}
	return j.Status == config.CronStatusUpdating || j.Status == config.CronStatusValid
}

package model

import "github.com/sqc157400661/jobx/config"

// JobCron
type JobCron struct {
	ID          int64  `gorm:"primaryKey;column:id" json:"id" xorm:"id pk autoincr"`
	EntryID     int    `gorm:"column:entry_id" json:"entry_id" xorm:"entry_id"`             // 定时任务id
	Spec        string `gorm:"column:spec" json:"spec" xorm:"spec"`                         // 定时表达式
	ExecType    string `gorm:"column:exec_type" json:"exec_type" xorm:"exec_type"`          // 执行任务类型，如job、func、shell
	ExecContent string `gorm:"column:exec_content" json:"exec_content" xorm:"exec_content"` // 执行任务内容
	Status      string `gorm:"column:status" json:"status" xorm:"status"`                   // 状态
	AppName     string `gorm:"column:app_name" json:"app_name" xorm:"app_name"`             // 应用名称
	Tenant      string `gorm:"column:tenant" json:"tenant" xorm:"tenant"`                   // tenant
	Locker      string `gorm:"column:locker" json:"locker" xorm:"locker"`                   // 锁拥有者
	CreateAt    int    `gorm:"column:create_at" json:"create_at" xorm:"created"`            // 创建时间
	UpdateAt    int    `gorm:"column:update_at" json:"update_at" xorm:"updated"`            // 更新时间
}

func (j *JobCron) TableName() string {
	return "job_cron"
}

func (j *JobCron) TouchEntryID(id int) (err error) {
	if j == nil {
		return
	}
	j.EntryID = id
	_, err = JFDb.ID(j.ID).Update(j)
	return
}

func (j *JobCron) MarkRunning() (err error) {
	if j == nil {
		return
	}
	j.Status = config.CronStatusRunning
	_, err = JFDb.ID(j.ID).Update(j)
	return
}

func (j *JobCron) MarkValid() (err error) {
	if j == nil {
		return
	}
	j.Locker = ""
	j.Status = config.CronStatusValid
	_, err = JFDb.ID(j.ID).Update(j)
	return
}

func (j *JobCron) MarkUpdating() (err error) {
	if j == nil {
		return
	}
	j.Status = config.CronStatusUpdating
	_, err = JFDb.ID(j.ID).Update(j)
	return
}

func (j *JobCron) MarkDeleted() (err error) {
	if j == nil {
		return
	}
	j.Status = config.CronStatusDeleted
	_, err = JFDb.ID(j.ID).Update(j)
	return
}

func (j *JobCron) MarkRebooting() (err error) {
	if j == nil {
		return
	}
	j.Status = config.CronStatusRebooting
	_, err = JFDb.ID(j.ID).Update(j)
	return
}

// Add 添加cron任务
func (j *JobCron) Add() (id int64, err error) {
	if j == nil {
		return
	}
	id, err = JFDb.Insert(j)
	return
}

// Update 更新
func (j *JobCron) Update() (err error) {
	if j == nil {
		return
	}
	_, err = JFDb.ID(j.ID).Update(j)
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

// IsUpdating 是否在更新状态
func (j *JobCron) IsUpdating() bool {
	if j == nil {
		return false
	}
	return j.Status == config.CronStatusUpdating
}

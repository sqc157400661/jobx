package dao

import "github.com/sqc157400661/jobx/config"

// Job [...]
type Job struct {
	ID       int    `gorm:"primaryKey;column:id" json:"id" xorm:"id pk autoincr"`
	RootID   int    `gorm:"column:root_id" json:"root_id" xorm:"root_id"`        // 根任务
	ParentID int    `gorm:"column:parent_id" json:"parent_id" xorm:"parent_id"`  // 父级任务
	Runnable int8   `gorm:"column:runnable" json:"runable" xorm:"runnable"`      // 是否是用于pipeline的任务
	Name     string `gorm:"column:name" json:"name" xorm:"name"`                 // 任务名称
	Desc     string `gorm:"column:description" json:"desc" xorm:"description"`   // 任务描述
	Owner    string `gorm:"column:owner" json:"owner" xorm:"owner"`              // 任务归属人
	AppName  string `gorm:"column:app_name" json:"app_name" xorm:"app_name"`     // 应用名称
	Tenant   string `gorm:"column:tenant" json:"tenant" xorm:"tenant"`           // tenant
	BizID    string `gorm:"column:biz_id" json:"biz_id,omitempty" xorm:"biz_id"` // 业务产生的唯一id
	Pause    int8   `gorm:"column:pause" json:"pause" xorm:"pause"`              // 是否允许暂停
	Locker   string `gorm:"column:locker" json:"locker" xorm:"locker"`           // 锁拥有者
	Level    int    `xorm:"level"`
	Path     string `xorm:"path"`
	State    `xorm:"extends"`
	Input    map[string]interface{} `gorm:"column:input" json:"input" xorm:"input"`           // 入参
	Env      map[string]interface{} `gorm:"column:env" json:"env" xorm:"env"`                 // 配置信息
	CreateAt int                    `gorm:"column:create_at" json:"create_at" xorm:"created"` // 创建时间
	UpdateAt int                    `gorm:"column:update_at" json:"update_at" xorm:"updated"` // 更新时间
}

func (j *Job) TableName() string {
	return "job"
}

func (j *Job) Save() (err error) {
	_, err = JFDb.InsertOne(j)
	return
}

func (t *Job) Update() (err error) {
	_, err = JFDb.Update(t, &Job{ID: t.ID})
	return
}

func (j *Job) MarkRunning() (err error) {
	j.State.Phase = config.PhaseRunning
	_, err = JFDb.Update(j, &Job{ID: j.ID, State: State{Phase: config.PhaseReady}})
	return
}

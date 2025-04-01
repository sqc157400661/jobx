package model

import (
	"github.com/sqc157400661/jobx/config"
	"k8s.io/klog/v2"
)

// PipelineTask [...]
type PipelineTask struct {
	ID       int    `gorm:"primaryKey;column:id" json:"id" xorm:"id pk autoincr"`
	JobID    int    `gorm:"column:job_id" json:"job_id" xorm:"job_id"`         // 任务id
	Name     string `gorm:"column:name" json:"name" xorm:"name"`               // 执行的action标志
	Action   string `gorm:"column:action" json:"action" xorm:"action"`         // 任务名称
	Desc     string `gorm:"column:description" json:"desc" xorm:"description"` // 任务描述
	Pause    int8   `gorm:"column:pause" json:"pause" xorm:"pause"`            // 是否允许暂停
	Retry    int    `gorm:"column:retry" json:"retry" xorm:"retry"`            // 允许自动重试次数
	Retries  int8   `gorm:"column:retries" json:"retries" xorm:"retries"`      // 已经自动重试的次数
	State    `xorm:"extends"`
	Input    map[string]interface{} `gorm:"column:input" json:"input" xorm:"input"`           // 入参
	Output   map[string]interface{} `gorm:"column:output" json:"output" xorm:"output"`        // 出参
	Env      map[string]interface{} `gorm:"column:env" json:"env" xorm:"env"`                 // 配置信息
	Context  map[string]interface{} `gorm:"column:context" json:"context" xorm:"context"`     // 上下文参数
	CreateAt int                    `gorm:"column:create_at" json:"create_at" xorm:"created"` // 创建时间
	UpdateAt int                    `gorm:"column:update_at" json:"update_at" xorm:"updated"` // 更新时间
}

func (t *PipelineTask) TableName() string {
	return "job_task"
}

//func (t *PipelineTask) MarkRunning() (err error) {
//	t.State.Phase = config.PhaseRunning
//	_, err = JFDb.Update(t, &PipelineTask{ID: t.ID, State: State{Phase: config.PhaseReady}})
//	return
//}

func (t *PipelineTask) Update() (err error) {
	_, err = JFDb.Update(t, &PipelineTask{ID: t.ID})
	return
}

func (t *PipelineTask) Next() (next *PipelineTask, err error) {
	next = new(PipelineTask)
	_, err = JFDb.Where("id>? and job_id=?", t.ID, t.JobID).Asc("id").Get(next)
	return
}

// Succeed PipelineTask has been successful
func (t *PipelineTask) Succeed() (err error) {
	t.State = State{
		Phase:  config.PhaseTerminated,
		Status: config.StatusSuccess,
	}
	_, err = JFDb.Update(t, &PipelineTask{ID: t.ID})
	return
}

func (t *PipelineTask) CanPause() bool {
	return t.Pause > 0
}

// Fail PipelineTask has failed
func (t *PipelineTask) Fail(err error) {
	var reason string
	if err != nil {
		reason = err.Error()
	}
	// 防止字符串超限
	if len(reason) >= 2000 {
		reason = reason[:2000]
	}
	t.State = State{
		Phase:  config.PhaseRunning,
		Status: config.StatusFail,
		Reason: reason,
	}
	_, e := JFDb.Update(t, &PipelineTask{ID: t.ID})
	if e != nil {
		klog.Errorf("PipelineTask Fail err：%s  id:%d ", e.Error(), t.ID)
	}
}

func (t *PipelineTask) RefreshState() (state State) {
	_, _ = JFDb.ID(t.ID).Get(t)
	return t.State
}

func (t *PipelineTask) Save() (err error) {
	_, err = JFDb.InsertOne(t)
	return
}

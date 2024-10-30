package dao

// JobLogs [...]
type JobLogs struct {
	ID       int    `xorm:"id pk autoincr" json:"-"`
	EventID  int    `xorm:"event_id" json:"event_id"` // 任务id
	Result   string `xorm:"result" json:"result"`     // 结果
	CreateAt int    `xorm:"created" json:"create_at"` // 创建时间
}

func (j *JobLogs) TableName() string {
	return "job_logs"
}

// Add 添加log
func (j *JobLogs) Add() (id int64, err error) {
	id, err = JFDb.Insert(j)
	return
}

// Update 更新log
func (j *JobLogs) Update() (err error) {
	_, err = JFDb.ID(j.ID).Update(j)
	return
}

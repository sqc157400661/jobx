package model

// JobToken
type JobToken struct {
	ID       int    `gorm:"primaryKey;column:id" json:"id" xorm:"id pk autoincr"`
	RootID   int    `gorm:"column:root_id" json:"root_id" xorm:"root_id"`     // 根任务
	Token    string `gorm:"column:token" json:"token" xorm:"token"`           // 令牌
	CreateAt int    `gorm:"column:create_at" json:"create_at" xorm:"created"` // 创建时间
	UpdateAt int    `gorm:"column:update_at" json:"update_at" xorm:"updated"` // 更新时间
}

func (j *JobToken) TableName() string {
	return "job_token"
}

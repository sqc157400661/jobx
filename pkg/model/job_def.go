package model

// JobDefinition [...]
type JobDefinition struct {
	ID       int    `xorm:"id pk autoincr" json:"id"`
	Name     string `xorm:"name" json:"name"`                         // 标识任务类型
	AppName  string `xorm:"app_name" json:"app_name" xorm:"app_name"` // 应用名称
	Tenant   string `xorm:"tenant" json:"tenant"`                     // 租户
	Version  int    `xorm:"version" json:"version"`                   // 版本
	YamlConf string `xorm:"yaml_conf" json:"yaml_conf"`               // 任务定义
}

func (j *JobDefinition) TableName() string {
	return "job_definition"
}

func (j *JobDefinition) Save() (err error) {
	_, err = DB().InsertOne(j)
	return
}

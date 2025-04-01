package model

// JobDefinition [...]
type JobDefinition struct {
	ID        int                    `xorm:"id pk autoincr" json:"id"`
	Name      string                 `xorm:"name" json:"name"`                         // 标识任务类型
	AppName   string                 `xorm:"app_name" json:"app_name" xorm:"app_name"` // 应用名称
	Tenant    string                 `xorm:"tenant" json:"tenant"`                     // 租户
	Sort      int                    `xorm:"sort" json:"sort"`                         // 排序
	Retry     int                    `xorm:"retry" json:"retry"`                       // 预设重试次数
	Pipelines []string               `xorm:"pipelines" json:"pipelines"`               // 任务流定义["k:v"]
	Input     map[string]interface{} `json:"input" xorm:"input"`                       // 预设参数
	Env       map[string]interface{} `xorm:"env" json:"env"`                           // 附加的env参数
	Condition map[string]string      `xorm:"condition" json:"condition"`               // 处理特殊逻辑
}

func (j *JobDefinition) TableName() string {
	return "job_definition"
}

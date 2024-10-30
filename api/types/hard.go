package types

type HardJobDefinition struct {
	Name      string            `json:"name"`   // 标识任务类型
	Tenant    string            `json:"tenant"` // 租户
	Pipelines []HardJobPipeline `json:"pipelines"`
}

type HardJobPipeline struct {
	Name      string                 `json:"name"`
	Action    string                 `json:"action"`
	Input     map[string]interface{} `json:"predefined_input"`
	Retry     int                    `json:"retry"`     // 预设重试次数
	Env       map[string]interface{} `json:"env"`       // 附加的env参数
	Condition map[string]string      `json:"condition"` // 处理特殊逻辑
}

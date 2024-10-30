package types

type CreateOptions struct {
	Name  string                 `json:"name"`
	Desc  string                 `json:"desc"`
	Pause bool                   `json:"pause"`
	Input map[string]interface{} `json:"input"`
	Env   map[string]interface{} `json:"env"`
}

type CreateJober struct {
	Owner  string `json:"owner"`
	Tenant string `json:"tenant"`
	BizId  string `json:"biz_id"`
	CreateOptions
	ChildJobs []*CreateJober
	Pipelines []*CreatePipeline
}

type CreatePipeline struct {
	Action string `json:"action"`
	Retry  int    `json:"retry"`
	CreateOptions
}

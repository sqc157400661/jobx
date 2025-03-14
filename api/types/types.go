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

type CommonParam struct {
	Timestamp string `json:"timestamp"`
	UID       string `json:"UID"`
	Action    string `json:"Action,omitempty"`
	Method    string `json:"method"`
}

type Pagination struct {
	PageIndex int `form:"pageIndex" json:"pageIndex" query:"pageIndex"`
	PageSize  int `form:"pageSize"  json:"pageSize"  query:"pageIndex"`
}

func (m *Pagination) GetPageIndex() int {
	if m.PageIndex <= 0 {
		m.PageIndex = 1
	}
	return m.PageIndex
}

func (m *Pagination) GetPageSize() int {
	if m.PageSize <= 0 {
		m.PageSize = 10
	}
	return m.PageSize
}

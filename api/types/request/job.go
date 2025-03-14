package request

import (
	"github.com/sqc157400661/jobx/api/types"
)

type GetJobListReq struct {
	types.Pagination `search:"-"`
	ParentId         int64  `query:"parent_id" json:"parent_id"`
	ID               string `query:"id" json:"id" form:"id"`
	Name             string `query:"name" json:"name" form:"name"`
	// operator, status
	Owner  string `query:"owner" json:"owner" form:"owner"`
	Status string `query:"status" json:"status" form:"status"`
	// input
	InputContain string `query:"input_contain" json:"input_contain" form:"input_contain"`
	Tenant       string `query:"tenant" json:"tenant" form:"tenant"`
	MinCreateAt  string `query:"min_create_at" json:"min_create_at" form:"min_create_at"` // min创建时间
	MaxCreateAt  string `query:"max_create_at" json:"max_create_at" form:"max_create_at"` // max创建时间
}

type GetJobReq struct {
	// ID Deprecated
	ID    int `query:"task_id" json:"task_id" form:"task_id"`
	JobID int `query:"id" json:"id" form:"id"`
}

type JobName struct {
	Code     string `json:"code"`
	Describe string `json:"desc"`
}

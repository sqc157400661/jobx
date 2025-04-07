package request

import (
	"github.com/sqc157400661/jobx/api/types"
)

type GetJobListReq struct {
	types.Pagination `search:"-"`
	ParentId         int64  `query:"parent_id" json:"parent_id"`
	IDs              []int  `query:"ids" json:"ids" search:"type:in;column:id;table:job"`
	Name             string `query:"name" json:"name" search:"type:contains;column:name;table:job"`
	Owner            string `query:"owner" json:"owner" search:"type:exact;column:owner;table:job"`
	Status           string `query:"status" json:"status" search:"type:exact;column:status;table:job"`
	InputContain     string `query:"input_contain" json:"input_contain" search:"type:contains;column:input;table:job"`
	AppName          string `query:"app_name" json:"app_name" search:"type:exact;column:app_name;table:job"`
	Tenant           string `query:"tenant" json:"tenant" search:"type:exact;column:tenant;table:job"`
	MinCreateAt      int    `query:"min_create_at" json:"min_create_at" search:"type:gte;column:create_at;table:job"` // min创建时间
	MaxCreateAt      int    `query:"max_create_at" json:"max_create_at" search:"type:lte;column:create_at;table:job"` // max创建时间
}

type GetJobReq struct {
	ID int `query:"id" json:"id" search:"type:exact;column:id;table:job"`
}

type RetryReq struct {
	TaskID int                    `query:"task_id" json:"task_id" form:"task_id"`
	Input  map[string]interface{} `query:"input" json:"input" form:"input"`
}

type SkipReq struct {
	TaskID int `query:"task_id" json:"task_id" form:"task_id"`
}

type PauseReq struct {
	JobID  int `query:"job_id" json:"job_id" form:"job_id"`
	TaskID int `query:"task_id" json:"task_id" form:"task_id"`
}

type RestartReq struct {
	JobID  int `query:"job_id" json:"job_id" form:"job_id"`
	TaskID int `query:"task_id" json:"task_id" form:"task_id"`
}

type DiscardReq struct {
	JobID  int `query:"job_id" json:"job_id" form:"job_id"`
	TaskID int `query:"task_id" json:"task_id" form:"task_id"`
}

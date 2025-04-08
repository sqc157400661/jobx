package request

import "github.com/sqc157400661/jobx/api/types"

type GetCronListReq struct {
	types.Pagination `search:"-"`
	AppName          string `query:"app_name" json:"app_name" search:"type:exact;column:app_name;table:job_cron"`
	Tenant           string `query:"tenant" json:"tenant" search:"type:exact;column:tenant;table:job_cron"`
	Name             string `json:"name" form:"name" query:"name" search:"type:exact;column:name;table:job_cron"`
	Owner            string `json:"owner" form:"owner" query:"owner" search:"type:exact;column:owner;table:job_cron"`
	Status           string `json:"status" form:"status" query:"status" search:"type:exact;column:status;table:job_cron"`
	ExecType         string `json:"exec_type" form:"exec_type" query:"exec_type" search:"type:exact;column:exec_type;table:job_cron"`
	ExecContent      string `json:"exec_content" form:"exec_content" query:"exec_content" search:"type:exact;column:exec_content;table:job_cron"`
}

type AddPlanReq struct {
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Spec        string `json:"spec"`         // 定时表达式
	ExecType    string `json:"exec_type"`    // 执行任务类型，如job、func、shell
	ExecContent string `json:"exec_content"` // 执行任务内容
	Status      string `json:"status"`       // 状态
	AppName     string `json:"app_name"`     // 应用名称
	Tenant      string `json:"tenant"`       // tenant
}

package types

//type JobReq struct {
//	ID    int64  `query:"id" json:"id"`         // 任务ID
//	BizID string `query:"biz_id" json:"biz_id"` // 业务产生的唯一id
//}

type JobListReq struct {
	ParentId int64  `query:"parent_id" json:"parent_id" form:"parent_id"`
	Page     int    `query:"currentPage" json:"currentPage" form:"currentPage"`
	Size     int    `query:"pageSize" json:"pageSize" form:"pageSize"`
	ID       int64  `query:"id" json:"id" form:"id"`
	Name     string `query:"name" json:"name" form:"name"`
	// Job ID list
	IDs []int64 `query:"IDs" json:"IDs" form:"IDs"`
	// operator, status
	Owner  string `query:"owner" json:"owner" form:"owner"`
	Status string `query:"status" json:"status" form:"status"`
	// input
	InputContain string `query:"input_contain" json:"input_contain" form:"input_contain"`
	Tenant       string `query:"tenant" json:"tenant" form:"tenant"`
	MinCreateAt  string `query:"min_create_at" json:"min_create_at" form:"min_create_at"` // min创建时间
	MaxCreateAt  string `query:"max_create_at" json:"max_create_at" form:"max_create_at"` // max创建时间
}

type TaskListReq struct {
	JobId int64 `query:"job_id" json:"job_id" form:"job_id"`
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

type LogListReq struct {
	Page   int `query:"currentPage" json:"currentPage" form:"currentPage"`
	Size   int `query:"pageSize" json:"pageSize" form:"pageSize"`
	RootID int `query:"root_id" json:"root_id" form:"root_id"`
}

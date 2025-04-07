package request

type TaskListReq struct {
	JobId int64 `query:"job_id" json:"job_id" search:"type:exact;column:job_id;table:job_task"`
}

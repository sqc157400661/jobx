package request

type GetJobLogReq struct {
	EventID int `json:"event_id" query:"event_id" search:"type:exact;column:event_id;table:job_logs"`
}

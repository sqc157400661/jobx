package request

type GetJobLogReq struct {
	EventID int `json:"event_id" form:"event_id" query:"event_id"`
}

package biz

import (
	"github.com/go-xorm/xorm"
	"github.com/sqc157400661/jobx/api/types"
	"github.com/sqc157400661/jobx/pkg/model"
)

func LogList(req types.LogListReq) (logs []model.JobLogs, total int64, err error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 20
	}
	engine := model.JFDb
	var sess *xorm.Session
	sess = engine.NewSession()
	if req.RootID > 0 {
		sess = engine.Where("event_id=?", req.RootID)
	}
	log := new(model.JobLogs)
	total, err = sess.Clone().Count(log)
	err = sess.Desc("id").Limit(req.Size, req.Size*(req.Page-1)).Find(&logs)
	return
}

func GetLogByEventID(eventID int) (jobLogs model.JobLogs, err error) {
	return model.GetLogByEventID(eventID)
}

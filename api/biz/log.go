package biz

import (
	"github.com/go-xorm/xorm"
	"github.com/sqc157400661/jobx/api/types"
	"github.com/sqc157400661/jobx/pkg/dao"
)

func LogList(req types.LogListReq) (logs []dao.JobLogs, total int64, err error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 20
	}
	model := dao.JFDb
	var sess *xorm.Session
	sess = model.NewSession()
	if req.RootID > 0 {
		sess = model.Where("event_id=?", req.RootID)
	}
	log := new(dao.JobLogs)
	total, err = sess.Clone().Count(log)
	err = sess.Desc("id").Limit(req.Size, req.Size*(req.Page-1)).Find(&logs)
	return
}

func GetLogByEventID(eventID int) (jobLogs dao.JobLogs, err error) {
	return dao.GetLogByEventID(eventID)
}

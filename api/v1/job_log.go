package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sqc157400661/jobx/api/base"
	"github.com/sqc157400661/jobx/api/biz"
	"github.com/sqc157400661/jobx/api/types/request"
	"github.com/sqc157400661/jobx/pkg/model"
)

type JobLog struct {
	base.Api
}

func (e JobLog) Get(c *gin.Context) {
	req := request.GetJobLogReq{}
	err := e.MakeContext(c).
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var jobLog model.JobLogs
	jobLog, err = biz.GetLogByEventID(req.EventID)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(jobLog, "查询成功")
}

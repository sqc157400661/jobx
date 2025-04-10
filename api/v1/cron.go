package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sqc157400661/jobx/api/base"
	"github.com/sqc157400661/jobx/api/types/request"
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/pkg/model"
	"github.com/sqc157400661/jobx/pkg/search"
)

type Cron struct {
	base.Api
}

func (e Cron) GetPage(c *gin.Context) {
	req := request.GetCronListReq{}
	err := e.MakeContext(c).
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	session := mysql.DB().NewSession()
	defer session.Close()
	var count int64
	list := make([]model.JobCron, 0)
	searchFunc := search.MakeCondition(req)
	count, err = searchFunc(session).Clone().Count(new(model.JobCron))
	if err != nil {
		e.Logger.Errorf("db error: %s", err)
		e.Error(500, err, "查询count失败")
		return
	}
	err = searchFunc(session).OrderBy("id desc").Find(&list)
	if err != nil {
		e.Logger.Errorf("db error: %s", err)
		e.Error(500, err, "查询失败")
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e Cron) Add(c *gin.Context) {
	req := request.AddPlanReq{}
	err := e.MakeContext(c).
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	if req.ExecType != config.JobExecType {
		err = fmt.Errorf("not support exec type: %s", req.ExecType)
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	// todo 查询任务def 判断ExecContent中的任务是否存在
	session := mysql.DB().NewSession()
	defer session.Close()
	var cron = &model.JobCron{
		Name:        req.Name,
		Spec:        req.Spec,
		Owner:       req.Owner,
		ExecType:    req.ExecType,
		ExecContent: req.ExecContent,
		AppName:     req.AppName,
		Tenant:      req.Tenant,
	}
	err = cron.Save()
	if err != nil {
		e.Logger.Errorf("db error: %s", err.Error())
		e.Error(500, err, "保存失败")
		return
	}
	e.OK(cron, "添加成功")
}

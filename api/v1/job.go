package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sqc157400661/jobx/api/base"
	"github.com/sqc157400661/jobx/api/biz"
	"github.com/sqc157400661/jobx/api/types/request"
	"github.com/sqc157400661/jobx/pkg/model"
	"github.com/sqc157400661/jobx/pkg/search"
)

type Job struct {
	base.Api
}

func (e Job) Get(c *gin.Context) {
	req := request.GetJobReq{}
	err := e.MakeContext(c).
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var jobRes model.Job
	session := mysql.DB().NewSession()
	defer session.Close()
	searchFunc := search.MakeCondition(req)
	_, err = searchFunc(session).Get(&jobRes)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(jobRes, "查询成功")
}

func (e Job) GetPage(c *gin.Context) {
	req := request.GetJobListReq{}
	err := e.MakeContext(c).
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var count int64
	var list []model.Job
	session := mysql.DB().NewSession()
	defer session.Close()
	searchFunc := search.MakeCondition(req, &req.Pagination)
	count, err = searchFunc(session).Clone().Count(new(model.Job))
	if err != nil {
		e.Logger.Errorf("db error: %s", err)
		e.Error(500, err, "查询count失败")
		return
	}
	err = searchFunc(session).Find(&list)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e Job) Retry(c *gin.Context) {
	req := request.RetryReq{}
	err := e.MakeContext(c).
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, nil, err.Error())
		return
	}
	err = biz.Retry(req)
	if err != nil {
		e.Error(500, nil, err.Error())
		return
	}
	e.OK(nil, "操作成功")
}

func (e Job) Skip(c *gin.Context) {
	req := request.SkipReq{}
	err := e.MakeContext(c).
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, nil, err.Error())
		return
	}
	err = biz.Skip(req)
	if err != nil {
		e.Error(500, nil, err.Error())
		return
	}
	e.OK(nil, "操作成功")
}

func (e Job) Discard(c *gin.Context) {
	req := request.DiscardReq{}
	err := e.MakeContext(c).
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, nil, err.Error())
		return
	}
	err = biz.Discard(req)
	if err != nil {
		e.Error(500, nil, err.Error())
		return
	}
	e.OK(nil, "操作成功")
}

func (e Job) Pause(c *gin.Context) {
	req := request.PauseReq{}
	err := e.MakeContext(c).
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, nil, err.Error())
		return
	}
	err = biz.PauseJob(req)
	if err != nil {
		e.Error(500, nil, err.Error())
		return
	}
	e.OK(nil, "操作成功")
}

func (e Job) Restart(c *gin.Context) {
	req := request.RestartReq{}
	err := e.MakeContext(c).
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, nil, err.Error())
		return
	}
	err = biz.RestartJob(req)
	if err != nil {
		e.Error(500, nil, err.Error())
		return
	}
	e.OK(nil, "操作成功")
}

func (e Job) ForceDiscard(c *gin.Context) {
	req := request.DiscardReq{}
	err := e.MakeContext(c).
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	err = biz.ForceDiscard(req)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	e.OK(nil, "操作成功")
}

func (e Job) TaskList(c *gin.Context) {
	req := request.TaskListReq{}
	err := e.MakeContext(c).
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var list []model.PipelineTask
	session := mysql.DB().NewSession()
	defer session.Close()
	searchFunc := search.MakeCondition(req)
	err = searchFunc(session).Find(&list)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(list, "查询成功")
}

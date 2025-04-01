package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sqc157400661/jobx/api/base"
	"github.com/sqc157400661/jobx/api/biz"
	"github.com/sqc157400661/jobx/api/types"
	"github.com/sqc157400661/jobx/api/types/request"
	"github.com/sqc157400661/jobx/pkg/model"
	"strconv"
	"strings"
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
	if req.ID == -1 {
		e.OK(map[string]string{
			"status":  "success",
			"message": "success",
			"tips":    "任务id为-1，默认返回成功，标志任务相关任务已经在执行且已经成功，又重新发起相同任务",
		}, "查询成功")
		return
	}
	if req.ID == 0 {
		req.ID = req.JobID
		if req.JobID == 0 {
			e.Error(400, nil, "任务参数不合法")
			return
		}
	}
	var jobRes types.JobResult
	jobRes, err = biz.Get(req.ID)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.OK(map[string]string{
		"status":  jobRes.Job.Status,
		"message": jobRes.Job.Reason,
		"tips":    "",
	}, "查询成功")
}

func (e Job) GetPage(c *gin.Context) {
	req := request.GetJobListReq{}
	err := e.MakeContext(c).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var count int64
	var list []model.Job
	var ids []int64
	if req.ID != "" {
		idArr := strings.Split(req.ID, ",")
		for _, v := range idArr {
			id, errC := strconv.Atoi(strings.Trim(v, " "))
			if errC != nil {
				continue
			}
			ids = append(ids, int64(id))
		}
	}
	list, count, err = biz.JobList(types.JobListReq{
		Name:         req.Name,
		Page:         req.PageIndex,
		Size:         req.PageSize,
		ParentId:     req.ParentId,
		IDs:          ids,
		Owner:        req.Owner,
		Status:       req.Status,
		InputContain: req.InputContain,
		Tenant:       req.Tenant,
		MinCreateAt:  req.MinCreateAt,
		MaxCreateAt:  req.MaxCreateAt,
	})
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e Job) Retry(c *gin.Context) {
	req := types.RetryReq{}
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
	req := types.SkipReq{}
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
	req := types.DiscardReq{}
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
	req := types.PauseReq{}
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
	req := types.RestartReq{}
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
	req := types.DiscardReq{}
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
	req := types.TaskListReq{}
	err := e.MakeContext(c).
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var tasks []model.PipelineTask
	tasks, err = biz.TaskList(req)
	e.OK(tasks, "查询成功")
}

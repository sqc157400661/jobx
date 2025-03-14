package biz

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/sqc157400661/jobx/api/types"
	"github.com/sqc157400661/jobx/pkg/dao"
	"time"
)

func JobList(req types.JobListReq) (jobs []dao.Job, total int64, err error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 20
	}
	layout := "2006-01-02 15:04:05"
	var t time.Time
	model := dao.JFDb
	var sess *xorm.Session
	sess = model.Where("parent_id=?", req.ParentId)
	if req.ID > 0 {
		sess = sess.ID(req.ID)
	}
	if len(req.IDs) > 0 {
		sess = sess.In("id", req.IDs)
	}
	if req.Name != "" {
		sess = sess.And("name=?", req.Name)
	}
	if req.Owner != "" {
		sess = sess.And("owner=?", req.Owner)
	}
	if req.Tenant != "" {
		sess = sess.And("tenant=?", req.Tenant)
	}
	if req.Status != "" {
		sess = sess.And("status=?", req.Status)
	}
	if req.MinCreateAt != "" {
		t, err = time.Parse(layout, req.MinCreateAt)
		if err == nil {
			// 使用 Unix 函数将 Time 类型转换为时间戳
			timestamp := t.Unix()
			fmt.Println(timestamp)
			sess = sess.And("create_at>=?", timestamp)
		}
	}
	if req.MaxCreateAt != "" {
		t, err = time.Parse(layout, req.MaxCreateAt)
		if err == nil {
			// 使用 Unix 函数将 Time 类型转换为时间戳
			timestamp := t.Unix()
			sess = sess.And("create_at<?", timestamp)
		}
	}
	if req.InputContain != "" {
		sess = sess.And("input like ?", "%"+req.InputContain+"%")
	}
	job := new(dao.Job)
	total, err = sess.Clone().Count(job)
	err = sess.Desc("id").Limit(req.Size, req.Size*(req.Page-1)).Find(&jobs)
	return
}

func TaskList(req types.TaskListReq) (tasks []dao.PipelineTask, err error) {
	model := dao.JFDb
	var sess *xorm.Session
	sess = model.NewSession()
	if req.JobId > 0 {
		sess = model.Where("job_id=?", req.JobId)
	}
	err = sess.Find(&tasks)
	return
}

package biz

import (
	"fmt"

	"github.com/sqc157400661/jobx/api/types"
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/internal/helper"
	"github.com/sqc157400661/jobx/pkg/model"
)

// Retry Only the failed task can be retried
func Retry(req types.RetryReq) (err error) {
	var task model.PipelineTask
	task, err = getTaskByID(req.TaskID)
	if err != nil {
		return
	}
	if !task.IsFailed() {
		err = fmt.Errorf("task status is not failed, retry is not allowed")
		return
	}
	var rootJob model.Job
	var has bool
	rootJob, has, err = model.GetRootJobByJobId(task.JobID)
	if err != nil {
		return
	}
	if !has {
		err = fmt.Errorf("root job not found jobID:%d", task.JobID)
		return
	}
	if rootJob.IsDiscarded() {
		err = fmt.Errorf("root job is discarded, retry is not allowed")
		return
	}
	if !rootJob.IsFailed() {
		err = fmt.Errorf("root job is locking for run, plase wait")
		return
	}
	state := model.State{
		Phase:  config.PhaseReady,
		Status: config.StatusPending,
	}
	task.State = state
	if req.Input != nil && len(req.Input) > 0 {
		task.Input = req.Input
	}
	_, err = model.DB().Cols("phase", "status", "input").Update(&task, &model.PipelineTask{ID: task.ID})
	if err != nil {
		return
	}
	sql := fmt.Sprintf("update %s set phase = ?,status = ?,locker='' where id in (?,?)", rootJob.TableName())
	_, err = model.DB().Exec(sql, state.Phase, state.Status, rootJob.ID, task.JobID)
	return
}

// Pause Only tasks can be pause
func Pause(req types.PauseReq) (err error) {
	var task model.PipelineTask
	task, err = getTaskByID(req.TaskID)
	if err != nil {
		return
	}
	if !task.CanPause() {
		err = fmt.Errorf("task can not pause taskID:%d", req.TaskID)
		return
	}
	if !task.IsReady() && !task.IsRunning() {
		err = fmt.Errorf("this state does not allow pausing taskID:%d", req.TaskID)
		return
	}
	task.Status = config.StatusPause
	_, err = model.DB().Cols("phase", "status").Update(&task, &model.PipelineTask{ID: task.ID})
	if err != nil {
		return
	}
	return
}

// PauseJob todo check root?
func PauseJob(req types.PauseReq) (err error) {
	var job model.Job
	var hasJob bool
	job, hasJob, err = model.GetJobById(req.JobID)
	if !hasJob || err != nil {
		return
	}
	if !job.IsRunning() && !job.IsReady() {
		err = fmt.Errorf("job(%d) status is not running, pause is not allowed", req.JobID)
		return
	}
	var tasks []*model.PipelineTask
	tasks, err = model.GetPipelineTasksByJobId(req.JobID)
	if err != nil {
		return
	}
	for _, task := range tasks {
		if !task.CanPause() {
			continue
		}
		if task.IsReady() {
			task.Status = config.StatusPause
			_, err = model.DB().Cols("phase", "status").Update(task, &model.PipelineTask{ID: task.ID})
			if err != nil {
				return
			}
		}
	}
	sql := fmt.Sprintf("update %s set status = ? where id =?", job.TableName())
	_, err = model.DB().Exec(sql, config.StatusPause, req.JobID)
	return
}

// RestartJob todo check root?
func RestartJob(req types.RestartReq) (err error) {
	var job model.Job
	var hasJob bool
	job, hasJob, err = model.GetJobById(req.JobID)
	if !hasJob || err != nil {
		return
	}
	if job.IsDiscarded() || job.IsFinished() {
		err = fmt.Errorf("job(%d) status is finished or discarded, restart is not allowed", req.JobID)
		return
	}
	var tasks []*model.PipelineTask
	tasks, err = model.GetPipelineTasksByJobId(req.JobID)
	if err != nil {
		return
	}
	for _, task := range tasks {
		if task.IsPausing() {
			task.Status = config.StatusPending
			task.Phase = config.PhaseReady
			_, err = model.DB().Cols("phase", "status").Update(task, &model.PipelineTask{ID: task.ID})
			if err != nil {
				return
			}
		}
	}
	sql := fmt.Sprintf("update %s set phase = ?,status = ?,locker='' where id=?", job.TableName())
	_, err = model.DB().Exec(sql, config.PhaseReady, config.StatusPending, req.JobID)
	if err != nil {
		return
	}
	return
}

// Skip skip one task
func Skip(req types.SkipReq) (err error) {
	var task model.PipelineTask
	var next *model.PipelineTask
	task, err = getTaskByID(req.TaskID)
	if err != nil {
		return
	}
	if task.IsSuccess() {
		err = fmt.Errorf("task status is success, skip is not allowed")
		return
	}
	if task.IsDiscarded() {
		err = fmt.Errorf("task status is discarded, skip is not allowed")
		return
	}
	var rootJob model.Job
	var has bool
	rootJob, has, err = model.GetRootJobByJobId(task.JobID)
	if err != nil {
		return
	}
	if !has {
		err = fmt.Errorf("root job not found jobID:%d", task.JobID)
		return
	}
	if rootJob.State.IsDiscarded() {
		err = fmt.Errorf("job status is discarded, skip is not allowed")
		return
	}
	task.Status = config.StatusSkip
	_, err = model.DB().Cols("phase", "status").Update(&task, &model.PipelineTask{ID: task.ID})
	if err != nil {
		return
	}
	next, err = task.Next()
	if err != nil {
		return
	}
	if next.ID > 0 {
		next.Context = helper.UnsafeMergeMap(next.Context, task.Context)
		_, err = model.DB().Cols("context").Update(next, &model.PipelineTask{ID: next.ID})
		if err != nil {
			return
		}
	}
	if rootJob.IsFailed() {
		state := model.State{
			Phase:  config.PhaseReady,
			Status: config.StatusPending,
		}
		sql := fmt.Sprintf("update %s set phase = ?,status = ?,locker='' where id in (?,?)", rootJob.TableName())
		_, err = model.DB().Exec(sql, state.Phase, state.Status, rootJob.ID, task.JobID)
	}
	return
}

// Discard Discard a job and clean up the bizID， Only the failed task can be discarded
func Discard(req types.DiscardReq) (err error) {
	var rootJob model.Job
	var has bool
	rootJob, has, err = model.GetRootJobByJobId(req.JobID)
	if err != nil {
		return
	}
	if !has {
		err = fmt.Errorf("root job not found jobID:%d", req.JobID)
		return
	}
	if !rootJob.IsFailed() {
		err = fmt.Errorf("root job is not failed, Abandon is not allowed")
		return
	}
	err = model.ReleaseTokens(rootJob.RootID)
	if err != nil {
		return
	}
	sql := fmt.Sprintf("update %s set status = ?,biz_id='' where id =?", rootJob.TableName())
	_, err = model.DB().Exec(sql, config.StatusDiscarded, rootJob.ID)
	return
}

// ForceDiscard Discard a job and clean up the bizID
func ForceDiscard(req types.DiscardReq) (err error) {
	var rootJob model.Job
	var has bool
	rootJob, has, err = model.GetRootJobByJobId(req.JobID)
	if err != nil {
		return
	}
	if !has {
		err = fmt.Errorf("root job not found jobID:%d", req.JobID)
		return
	}
	err = model.ReleaseTokens(rootJob.RootID)
	if err != nil {
		return
	}
	sql := fmt.Sprintf("update %s set status = ?,biz_id='' where id =?", rootJob.TableName())
	_, err = model.DB().Exec(sql, config.StatusDiscarded, rootJob.ID)
	if req.TaskID > 0 {
		_, err = model.DB().ID(req.TaskID).Cols("phase", "status").Update(map[string]string{
			"phase":  config.PhaseTerminated,
			"status": config.StatusDiscarded,
		})
	}
	if err != nil {
		return
	}
	return
}

// Discard  discard root Job
//func DiscardRootJobByID(RootJobID int) {
//
//}

func getTaskByID(id int) (task model.PipelineTask, err error) {
	_, err = model.DB().ID(id).Get(&task)
	if err != nil {
		return
	}
	if task.ID <= 0 {
		err = fmt.Errorf("not found task")
		return
	}
	return
}

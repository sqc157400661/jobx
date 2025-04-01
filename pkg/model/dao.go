package model

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/pkg/errors"
)

// todo 这里实例化多次jobflow会有值覆盖现象，但是目前单个项目对应单库，暂时可以不用考虑
var JFDb *xorm.Engine

// GetChildRunableJobsByRootId 根据rootId获取可执行的子任务
func GetChildRunableJobsByRootId(rootId int, status ...string) (jobs []Job, err error) {
	model := JFDb.Where("root_id = ? and runnable =1", rootId)
	if len(status) > 0 {
		model.In("phase", status)
	}
	err = model.Asc("id").Find(&jobs)
	return
}

func GetPipelineTasksByJobId(id int) (tasks []*PipelineTask, err error) {
	err = JFDb.Where("job_id=?", id).Asc("id").Find(&tasks)
	return
}

func GetJobById(id int) (job Job, has bool, err error) {
	has, err = JFDb.ID(id).Get(&job)
	return
}

func CheckTokens(tokens []string) (rootId int, err error) {
	var jobTokens []*JobToken
	err = JFDb.In("token", tokens).Find(&jobTokens)
	if err != nil {
		return
	}
	if len(jobTokens) > 0 {
		var jobIDs []int
		var jobIDsMap = map[int]bool{}
		var hasTokens []string
		var hasTokensMap = map[string]bool{}
		for _, v := range jobTokens {
			rootId = v.RootID
			if _, has := jobIDsMap[v.RootID]; !has {
				jobIDsMap[v.RootID] = true
				jobIDs = append(jobIDs, v.RootID)
			}
			if _, has := hasTokensMap[v.Token]; !has {
				hasTokensMap[v.Token] = true
				hasTokens = append(hasTokens, v.Token)
			}
		}
		err = errors.TokenConflict(fmt.Sprintf("相关任务：%v，冲突tokens：%v", jobIDs, hasTokens))
	}
	return
}

func CreateTokens(rootId int, tokens []string) (err error) {
	if len(tokens) == 0 {
		return
	}
	jobTokens := make([]*JobToken, len(tokens))
	for k, token := range tokens {
		jobTokens[k] = &JobToken{
			Token:  token,
			RootID: rootId,
		}
	}
	_, err = JFDb.Insert(&jobTokens)
	return
}

func ReleaseTokens(rootId int) (err error) {
	var jobToken = new(JobToken)
	_, err = JFDb.Where("root_id=?", rootId).Delete(jobToken)
	return
}

func GetRootJobByJobId(jobId int) (resJob Job, has bool, err error) {
	var job Job
	has, err = JFDb.ID(jobId).Get(&job)
	if job.RootID == 0 {
		resJob = job
		return
	}
	has, err = JFDb.ID(job.RootID).Get(&resJob)
	return
}

func GetJobByBizId(bizId string) (job Job, has bool, err error) {
	has, err = JFDb.Where("biz_id=?", bizId).Get(&job)
	return
}

func GetChildJobsById(id int) (childJobs []*Job, err error) {
	err = JFDb.Where("parent_id=?", id).Find(&childJobs)
	return
}

func UpdateJobStateByID(id int, state *State) (err error) {
	job := new(Job)
	job.State = *state
	_, err = JFDb.ID(id).Update(job)
	return
}

func UpdateJobsStateByRootID(rootid int, state *State) (err error) {
	_, err = JFDb.Table(new(Job)).Where("root_id=?", rootid).Update(map[string]interface{}{
		"phase":  state.Phase,
		"status": state.Status,
	})
	return
}

// GetValidCron 获取有效的cron任务
func GetValidCron() (jobCrons []*JobCron, err error) {
	model := JFDb.Where("status=？", config.CronStatusValid)
	err = model.Asc("id").Find(&jobCrons)
	return
}

// CollectCron 获取相关cron任务
func CollectCron(uid string) (jobCrons []*JobCron, err error) {
	// 尝试lock
	session := JFDb.NewSession()
	defer session.Close()
	_, err = session.Exec("update job_cron set locker=? where locker='' and status =? order by id asc limit ?", uid, config.CronStatusValid, 2)
	if err != nil {
		return
	}
	err = session.Commit()
	if err != nil {
		return
	}
	model := JFDb.Where("locker=?", uid)
	err = model.Asc("id").Find(&jobCrons)
	return
}

// GetCronByExecContent 根据content获取cron任务
func GetCronByExecContent(content string) (jobCron *JobCron, err error) {
	jobCron = &JobCron{}
	_, err = JFDb.Where("exec_content=?", content).Get(jobCron)
	return
}

// GetCronByID 根据id获取cron任务
func GetCronByID(id int) (jobCron JobCron, err error) {
	_, err = JFDb.ID(id).Get(&jobCron)
	return
}

// GetLogByEventID 根据EventID获取log信息
func GetLogByEventID(eventID int) (jobLogs JobLogs, err error) {
	_, err = JFDb.Where("event_id=?", eventID).Get(&jobLogs)
	return
}

// BatchAddLogs 批量添加日志
func BatchAddLogs(logs []*JobLogs) (err error) {
	_, err = JFDb.Insert(&logs)
	return
}

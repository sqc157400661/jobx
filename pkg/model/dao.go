package model

import (
	"fmt"
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/pkg/errors"
	"github.com/sqc157400661/jobx/pkg/mysql"
)

// GetChildRunableJobsByRootId 根据rootId获取可执行的子任务
func GetChildRunableJobsByRootId(rootId int, status ...string) (jobs []Job, err error) {
	model := mysql.DB().Where("root_id = ? and runnable =1", rootId)
	if len(status) > 0 {
		model.In("phase", status)
	}
	err = model.Asc("id").Find(&jobs)
	return
}

func GetPipelineTasksByJobId(id int) (tasks []*PipelineTask, err error) {
	err = mysql.DB().Where("job_id=?", id).Asc("id").Find(&tasks)
	return
}

func GetJobById(id int) (job Job, has bool, err error) {
	has, err = mysql.DB().ID(id).Get(&job)
	return
}

func CheckTokens(tokens []string) (rootId int, err error) {
	var jobTokens []*JobToken
	err = mysql.DB().In("token", tokens).Find(&jobTokens)
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
	_, err = mysql.DB().Insert(&jobTokens)
	return
}

func ReleaseTokens(rootId int) (err error) {
	var jobToken = new(JobToken)
	_, err = mysql.DB().Where("root_id=?", rootId).Delete(jobToken)
	return
}

func GetRootJobByJobId(jobId int) (resJob Job, has bool, err error) {
	var job Job
	has, err = mysql.DB().ID(jobId).Get(&job)
	if job.RootID == 0 {
		resJob = job
		return
	}
	has, err = mysql.DB().ID(job.RootID).Get(&resJob)
	return
}

func GetJobByBizId(bizId string) (job Job, has bool, err error) {
	has, err = mysql.DB().Where("biz_id=?", bizId).Get(&job)
	return
}

func GetChildJobsById(id int) (childJobs []*Job, err error) {
	err = mysql.DB().Where("parent_id=?", id).Find(&childJobs)
	return
}

func UpdateJobStateByID(id int, state *State) (err error) {
	job := new(Job)
	job.State = *state
	_, err = mysql.DB().ID(id).Update(job)
	return
}

func UpdateJobsStateByRootID(rootid int, state *State) (err error) {
	_, err = mysql.DB().Table(new(Job)).Where("root_id=?", rootid).Update(map[string]interface{}{
		"phase":  state.Phase,
		"status": state.Status,
	})
	return
}

// GetValidCron 获取有效的cron任务
func GetValidCron() (jobCrons []*JobCron, err error) {
	model := mysql.DB().Where("status=？", config.CronStatusValid)
	err = model.Asc("id").Find(&jobCrons)
	return
}

// CollectCron 获取相关cron任务
func CollectCron(uid string) (jobCrons []*JobCron, err error) {
	// 尝试lock
	session := mysql.DB().NewSession()
	defer session.Close()
	_, err = session.Exec("update job_cron set locker=? where locker='' and status =? order by id asc limit ?", uid, config.CronStatusValid, 2)
	if err != nil {
		return
	}
	err = session.Commit()
	if err != nil {
		return
	}
	model := mysql.DB().Where("locker=?", uid)
	err = model.Asc("id").Find(&jobCrons)
	return
}

// GetCronByExecContent 根据content获取cron任务
func GetCronByExecContent(content string) (jobCron *JobCron, err error) {
	jobCron = &JobCron{}
	_, err = mysql.DB().Where("exec_content=?", content).Get(jobCron)
	return
}

// GetCronByID 根据id获取cron任务
func GetCronByID(id int) (jobCron JobCron, err error) {
	_, err = mysql.DB().ID(id).Get(&jobCron)
	return
}

// GetLogByEventID 根据EventID获取log信息
func GetLogByEventID(eventID int) (jobLogs JobLogs, err error) {
	_, err = mysql.DB().Where("event_id=?", eventID).Get(&jobLogs)
	return
}

// BatchAddLogs 批量添加日志
func BatchAddLogs(logs []*JobLogs) (err error) {
	_, err = mysql.DB().Insert(&logs)
	return
}

// GetJobDefByName ...
func GetJobDefByName(name string) (jobCron *JobDefinition, err error) {
	jobCron = &JobDefinition{}
	_, err = mysql.DB().Where("name=?", name).Get(jobCron)
	return
}

// GetCronByName ...
func IsCronJobExist(tenant, appName, name string) (jobCron *JobCron, err error) {
	jobCron = &JobCron{}
	_, err = mysql.DB().Where("name=? and app_name=? and tenant=?", name, appName, tenant).Get(jobCron)
	return
}

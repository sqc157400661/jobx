package job

import (
	"fmt"
	"github.com/sqc157400661/jobx/cmd/service"
	"github.com/sqc157400661/jobx/hack/demo"
	"github.com/sqc157400661/jobx/pkg/mysql"
	"github.com/sqc157400661/jobx/pkg/options/jobopt"
	"github.com/sqc157400661/jobx/pkg/providers"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sqc157400661/jobx/test"
)

func TestDemoAutoSuspendJob(t *testing.T) {
	var err error
	input := jobopt.JobInput(map[string]interface{}{
		"action": "Suspend",
	})
	// test multiple pipeline add
	err = NewJob("AutoSuspend", "sqc", jobopt.JobDesc("CDWInternal"), input).
		AddPipeline("QueryIdlesMetric", "QueryIdlesMetric").
		AddPipeline("CheckIdle", "CheckIdle").
		AddPipeline("PreVwCheckTasker", "PreVwCheck").
		AddPipeline("LockVwStatus", "LockVwStatusInDB").
		AddPipeline("UpdateK8sResource", "UpdateK8sResource").
		AddPipeline("wait", "delay", jobopt.JobInput(map[string]interface{}{"time": 5 * time.Second})).
		AddPipeline("CheckK8sResource", "CheckK8sResource").
		AddPipeline("UpdateVwStatusInDB", "MarkVwStatusInDB").
		Exec()
	assert.NoError(t, err)
}

func TestDemoAutoResumeJob(t *testing.T) {
	var err error
	input := jobopt.JobInput(map[string]interface{}{
		"action": "Resume",
	})
	// test multiple pipeline add
	err = NewJob("AutoResume", "sqc", jobopt.JobDesc("CDWInternal"), input).
		AddPipeline("QueryCnchPendingTask", "QueryPendingTask").
		AddPipeline("PreVwCheckTasker", "PreVwCheck").
		AddPipeline("LockVwStatus", "LockVwStatusInDB").
		AddPipeline("UpdateK8sResource", "UpdateK8sResource").
		AddPipeline("wait", "delay", jobopt.JobInput(map[string]interface{}{"time": 5 * time.Second})).
		AddPipeline("CheckK8sResource", "CheckK8sResource").
		AddPipeline("UpdateVwStatusInDB", "MarkVwStatusInDB").
		Exec()
	assert.NoError(t, err)
}

func TestJober(t *testing.T) {
	var err error
	// test multiple pipeline add
	err = NewJob("jober1", "sqc").
		AddPipeline("task_1", "demo").
		AddPipeline("task_2", "delay").
		AddPipeline("task_3", "demo2").
		AddPipeline("task_4", "demo").
		Exec()
	assert.NoError(t, err)

	// test add jobInput and JobEnv
	inputMap := map[string]interface{}{
		"testKeyint":    2,
		"testKeybool":   true,
		"testKeystring": "hahah",
	}
	err = NewJob("jober2", "sqc",
		jobopt.JobInput(map[string]interface{}{
			"all_config": "yes",
		}),
		jobopt.JobEnv(map[string]interface{}{
			"env": "test",
		}),
	).
		AddJob(NewJob("jober2_1", "sqc", jobopt.JobInput(inputMap))).
		AddPipeline("task_1", "demo", jobopt.JobEnv(map[string]interface{}{
			"env": "sim",
		})).
		AddPipeline("task_2", "delay", jobopt.JobInput(map[string]interface{}{
			"pipeline": 233,
		})).
		AddPipeline("task_3", "demo").
		Exec()
	assert.NoError(t, err)

	// test multiple job add
	job := NewJob("jober3", "sqc")
	err = job.AddJob(
		NewJob("jober3_1", "sqc").
			AddPipeline("task_1", "demo").
			AddPipeline("task_2", "delay").
			AddPipeline("task_3", "demo2")).Exec()
	assert.NoError(t, err)
	err = job.AddJob(
		NewJob("jober3_2", "sqc").
			AddPipeline("task_1", "demo").
			AddPipeline("task_2", "demo").
			AddPipeline("task_3", "demo")).Exec()
	assert.NoError(t, err)

	// test job with biz_id
	err = NewJob(
		"jober1",
		"sqc",
		jobopt.BizId("12345"),
	).
		AddPipeline("task_1", "demo").Exec()
	assert.NoError(t, err)
}

func TestWaitJob(t *testing.T) {
	engine, err := test.GetEngine()
	require.NoError(t, err)
	mysql.JFDb = engine
	err = WaitJob(400039, 10)
	fmt.Println(err, 343454534)

}

func init() {
	engine, err := test.GetEngine()
	if err != nil {
		panic(err)
	}
	mysql.JFDb = engine
	jobFlow, _ := service.NewJobFlow("sqc_test_compute", engine)
	_ = jobFlow.Register(
		&providers.DemoTasker{},
		&providers.Demo2Tasker{},
		&demo.CheckIdle{},
		&demo.MarkVwStatusInDB{},
		&demo.QueryCnchPendingTask{},
		&demo.MarkVwPendingStatusInDB{},
		&demo.QueryMetric{},
		&demo.UpdateK8sResource{},
		&demo.UpdateK8sResourceCheckLoop{},
		&demo.PreVwCheckTasker{},
	)
}

func TestCronJob(t *testing.T) {
	input := jobopt.JobInput(map[string]interface{}{
		"action": "Suspend",
	})
	job := NewJob("AutoResume", "sqc", jobopt.JobDesc("CDWInternal"), input).
		AddPipeline("QueryCnchPendingTask", "QueryPendingTask").
		AddPipeline("PreVwCheckTasker", "PreVwCheck").
		AddPipeline("LockVwStatus", "LockVwStatusInDB").
		AddPipeline("UpdateK8sResource", "UpdateK8sResource").
		AddPipeline("wait", "delay", jobopt.JobInput(map[string]interface{}{"time": 5 * time.Second})).
		AddPipeline("CheckK8sResource", "CheckK8sResource").
		AddPipeline("UpdateVwStatusInDB", "MarkVwStatusInDB")
	cron, err := NewCronjob("* * * * * *", "testcron2", "sqc")
	fmt.Println(err)
	fmt.Println(cron.ExecJob(job))
}

func TestYAMLToJob(t *testing.T) {
	a := `jobDef:
    name: testcron3_AutoResume
    desc: CDWInternal
    owner: sqc
    appName: iam
    tenant: beijing
    bizId: ""
    pipelines:
        - name: QueryCnchPendingTask
          desc: 查询cnch的pengding-task表
          action: QueryPendingTask
          retryNum: 3
          input:
            action: Suspend
          env:
            retry_gap_second: 10
        - name: PreVwCheckTasker
          desc: 校验Vw的状态
          action: PreVwCheck
          retryNum: 3
          input:
            action: Suspend
          env:
            retry_gap_second: 10
        - name: LockVwStatus
          desc: lock vw status in db
          action: LockVwStatusInDB
          retryNum: 3
          input:
            action: Suspend
          env:
            retry_gap_second: 10
        - name: UpdateK8sResource
          desc: 更新对应的K8s VW资源信息
          action: UpdateK8sResource
          retryNum: 3
          input:
            action: Suspend
          env:
            retry_gap_second: 10
        - name: wait
          desc: 用于控制任务之间，需要等待的时间时长，返回结果里的wait_time，代表需要等待的时间
          action: delay
          retryNum: 3
          input:
            action: Suspend
            time: 5s
          env:
            retry_gap_second: 10
        - name: CheckK8sResource
          desc: watch k8s资源是否被更新完成
          action: CheckK8sResource
          retryNum: 3
          input:
            action: Suspend
          env:
            retry_gap_second: 10
        - name: UpdateVwStatusInDB
          desc: 更新vw在数据库中的状态信息
          action: MarkVwStatusInDB
          retryNum: 3
          input:
            action: Suspend
          env:
            retry_gap_second: 10
    input:
        action: Suspend
    env:
        retry_gap_second: 10

`
	fmt.Println(YAMLToJob([]byte(a)))
}

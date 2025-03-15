package job

import (
	"fmt"
	"github.com/sqc157400661/jobx/cmd/service"
	"github.com/sqc157400661/jobx/hack/demo"
	"github.com/sqc157400661/jobx/pkg/providers"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sqc157400661/jobx/pkg/dao"
	"github.com/sqc157400661/jobx/pkg/options"
	"github.com/sqc157400661/jobx/test"
)

func TestDemoAutoSuspendJob(t *testing.T) {
	var err error
	input := options.JobInput(map[string]interface{}{
		"action": "Suspend",
	})
	// test multiple pipeline add
	err = NewJober("AutoSuspend", "sqc", "CDWInternal", input).
		AddPipeline("QueryIdlesMetric", "QueryIdlesMetric").
		AddPipeline("CheckIdle", "CheckIdle").
		AddPipeline("PreVwCheckTasker", "PreVwCheck").
		AddPipeline("LockVwStatus", "LockVwStatusInDB").
		AddPipeline("UpdateK8sResource", "UpdateK8sResource").
		AddPipeline("wait", "delay", options.JobInput(map[string]interface{}{"time": 5 * time.Second})).
		AddPipeline("CheckK8sResource", "CheckK8sResource").
		AddPipeline("UpdateVwStatusInDB", "MarkVwStatusInDB").
		Exec()
	assert.NoError(t, err)
}

func TestDemoAutoResumeJob(t *testing.T) {
	var err error
	input := options.JobInput(map[string]interface{}{
		"action": "Resume",
	})
	// test multiple pipeline add
	err = NewJober("AutoSuspend", "sqc", "CDWInternal", input).
		AddPipeline("QueryCnchPendingTask", "QueryPendingTask").
		AddPipeline("PreVwCheckTasker", "PreVwCheck").
		AddPipeline("LockVwStatus", "LockVwStatusInDB").
		AddPipeline("UpdateK8sResource", "UpdateK8sResource").
		AddPipeline("wait", "delay", options.JobInput(map[string]interface{}{"time": 5 * time.Second})).
		AddPipeline("CheckK8sResource", "CheckK8sResource").
		AddPipeline("UpdateVwStatusInDB", "MarkVwStatusInDB").
		Exec()
	assert.NoError(t, err)
}

func TestJober(t *testing.T) {
	var err error
	// test multiple pipeline add
	err = NewJober("jober1", "sqc", "").
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
	err = NewJober("jober2", "sqc", "",
		options.JobInput(map[string]interface{}{
			"all_config": "yes",
		}),
		options.JobEnv(map[string]interface{}{
			"env": "test",
		}),
	).
		AddJob(NewJober("jober2_1", "sqc", "", options.JobInput(inputMap))).
		AddPipeline("task_1", "demo", options.JobEnv(map[string]interface{}{
			"env": "sim",
		})).
		AddPipeline("task_2", "delay", options.JobInput(map[string]interface{}{
			"pipeline": 233,
		})).
		AddPipeline("task_3", "demo").
		Exec()
	assert.NoError(t, err)

	// test multiple job add
	job := NewJober("jober3", "sqc", "")
	err = job.AddJob(
		NewJober("jober3_1", "sqc", "").
			AddPipeline("task_1", "demo").
			AddPipeline("task_2", "delay").
			AddPipeline("task_3", "demo2")).Exec()
	assert.NoError(t, err)
	err = job.AddJob(
		NewJober("jober3_2", "sqc", "").
			AddPipeline("task_1", "demo").
			AddPipeline("task_2", "demo").
			AddPipeline("task_3", "demo")).Exec()
	assert.NoError(t, err)

	// test job with biz_id
	err = NewJober(
		"jober1",
		"sqc",
		"",
		options.BizId("12345"),
	).
		AddPipeline("task_1", "demo").Exec()
	assert.NoError(t, err)
}

func TestWaitJob(t *testing.T) {
	engine, err := test.GetEngine()
	require.NoError(t, err)
	dao.JFDb = engine
	err = WaitJob(400039, 10)
	fmt.Println(err, 343454534)

}

func init() {
	engine, err := test.GetEngine()
	if err != nil {
		panic(err)
	}
	dao.JFDb = engine
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

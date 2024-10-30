package v1

import (
	"github.com/sqc157400661/jobx/api/types"
	"github.com/sqc157400661/jobx/pkg/dao"
	"github.com/sqc157400661/jobx/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJobList(t *testing.T) {
	engine, err := test.GetEngine()
	require.NoError(t, err)
	dao.JFDb = engine
	var req = types.JobListReq{
		Owner: "sqc",
		//IDs:   []int64{514, 510},
		MinCreateAt: "2023-04-21 15:07:46",
		MaxCreateAt: "2023-04-21 15:07:48",
	}
	jobs, tal, err := JobList(req)
	assert.NoError(t, err)
	t.Log(jobs, tal)
}

func TestTaskList(t *testing.T) {
	engine, err := test.GetEngine()
	require.NoError(t, err)
	dao.JFDb = engine
	var req = types.TaskListReq{
		JobId: 591,
	}
	tasks, err := TaskList(req)
	assert.NoError(t, err)
	t.Log(tasks)
}

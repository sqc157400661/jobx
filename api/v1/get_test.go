package v1

import (
	"github.com/sqc157400661/jobx/pkg/dao"
	"github.com/sqc157400661/jobx/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGet(t *testing.T) {
	engine, err := test.GetEngine()
	require.NoError(t, err)
	dao.JFDb = engine
	job, err := Get(514)
	assert.NoError(t, err)
	t.Log(job.Job)
	t.Log(job.ChildJobs)
	t.Log(job.Task)
	job, err = Get(515)
	assert.NoError(t, err)
	t.Log(job.Job)
	t.Log(job.ChildJobs)
	t.Log(job.Task)
}

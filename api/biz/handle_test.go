package biz

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sqc157400661/jobx/api/types"
	"github.com/sqc157400661/jobx/internal/helper"
	"github.com/sqc157400661/jobx/pkg/model"
	"github.com/sqc157400661/jobx/test"
)

func TestRetry(t *testing.T) {
	engine, err := test.GetEngine()
	require.NoError(t, err)
	model.JFDb = engine
	var req = types.RetryReq{
		TaskID: 935,
		//IDs:   []int64{514, 510},
	}
	err = Retry(req)
	assert.NoError(t, err)
}

func TestPause(t *testing.T) {
	engine, err := test.GetEngine()
	require.NoError(t, err)
	model.JFDb = engine
	var req = types.PauseReq{
		TaskID: 935,
		//IDs:   []int64{514, 510},
	}
	err = Pause(req)
	assert.NoError(t, err)
}

func TestPauseJob(t *testing.T) {
	engine, err := test.GetEngine()
	require.NoError(t, err)
	model.JFDb = engine
	var req = types.PauseReq{
		JobID: 132,
		//IDs:   []int64{514, 510},
	}
	err = PauseJob(req)
	assert.NoError(t, err)
}

func TestRestartJob(t *testing.T) {
	engine, err := test.GetEngine()
	require.NoError(t, err)
	model.JFDb = engine
	var req = types.RestartReq{
		JobID: 132,
		//IDs:   []int64{514, 510},
	}
	err = RestartJob(req)
	assert.NoError(t, err)
}

func TestAbandon(t *testing.T) {
	engine, err := test.GetEngine()
	require.NoError(t, err)
	model.JFDb = engine
	var req = types.DiscardReq{
		JobID: 11111193,
		//IDs:   []int64{514, 510},
	}
	err = Discard(req)
	assert.NoError(t, err)
}

func TestSkip(t *testing.T) {
	engine, err := test.GetEngine()
	require.NoError(t, err)
	now := new(model.PipelineTask)
	model.JFDb = engine
	_, err = model.DB().Where("id=?", 1443).Asc("id").Get(now)
	next, _ := now.Next()
	fmt.Println()
	next.Context = helper.UnsafeMergeMap(next.Context, now.Context)
	fmt.Println(next.Context)
	_, err = model.DB().Cols("context").ID(next.ID).Update(next)
	assert.NoError(t, err)
}

func TestA(t *testing.T) {
	str1 := "v3.2.4.5"
	str := "5.7.25-OceanBase-v3.2.4.5"
	length := len("-OceanBase-")
	index := strings.Index(str, "-OceanBase-")
	fmt.Println(strings.Index(str1, "-OceanBase-"))
	fmt.Println(str[index+length:])
	fmt.Println(time.Now().Format("200601021504"))

}

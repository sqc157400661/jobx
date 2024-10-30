package service

import (
	"github.com/sqc157400661/jobx/pkg/providers"
	"github.com/sqc157400661/jobx/test"
	"github.com/stretchr/testify/require"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestMutiJobFlow(t *testing.T) {
	engine, err := test.GetEngine()
	require.NoError(t, err)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			jobFlow, err := NewJobFlow("test_id_"+strconv.Itoa(i), engine)
			require.NoError(t, err)
			_ = jobFlow.AddProvider(&providers.DemoTasker{}, "demo")
			_ = jobFlow.AddProvider(&providers.Demo2Tasker{}, "demo2")
			require.NoError(t, err)
			jobFlow.Start()
			time.Sleep(30 * time.Second)
			t.Log(" jobFlow Quit")
			jobFlow.Quit()
		}(i)
	}
	wg.Wait()
}

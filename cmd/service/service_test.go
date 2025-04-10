package service

import (
	"github.com/sqc157400661/jobx/config"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/sqc157400661/jobx/pkg/providers"
)

func TestMutiJobFlow(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			jobFlow, err := NewJobFlow("test_id_"+strconv.Itoa(i), config.MySQL{
				Host:   "localhost",
				User:   "root",
				Passwd: "157400661",
				DB:     "task_center",
				Port:   3306,
			})
			require.NoError(t, err)
			_ = jobFlow.AddProvider(&providers.DemoTasker{}, "demo")
			_ = jobFlow.AddProvider(&providers.Demo2Tasker{}, "demo2")
			require.NoError(t, err)
			jobFlow.Start()
			time.Sleep(30 * time.Second)
			t.Log(" jobFlow Quit")
			jobFlow.Stop()
		}(i)
	}
	wg.Wait()
}

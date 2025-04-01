package cron

import (
	"fmt"
	"github.com/sqc157400661/jobx/pkg/model"
	"github.com/sqc157400661/jobx/test"
	"github.com/stretchr/testify/require"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	engine, err := test.GetEngine()
	require.NoError(t, err)
	model.JFDb = engine
	runner, err := NewCronRunner("1")
	require.NoError(t, err)
	err = runner.Start()
	require.NoError(t, err)
	runner.Add("1,2,3 * * * * 7,1,2,3,4,5", "test", &testJob{})
	time.Sleep(time.Minute * 10)
}

type testJob struct {
}

func (j *testJob) Run() {
	now := time.Now()
	fmt.Printf("second:%d test...\n", now.Second())
}

// GetSecondLevelDomain 提取给定域名的二级域名部分
func GetSecondLevelDomain(referer string) string {
	// 解析URL
	parsedURL, err := url.Parse(referer)
	if err != nil {
		return referer
	}
	// 提取主机名
	host := parsedURL.Host

	// 分割主机名获取二级域名
	parts := strings.Split(host, ".")
	if len(parts) < 3 {
		return referer
	}
	return "." + parts[len(parts)-2] + "." + parts[len(parts)-1]
}

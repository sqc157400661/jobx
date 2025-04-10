package log

import (
	"github.com/sqc157400661/jobx/pkg/mysql"
	"github.com/sqc157400661/jobx/test"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	engine, err := test.GetEngine()
	mysql.JFDb = engine
	require.NoError(t, err)
	Info(3333, "测试12345678")
	time.Sleep(time.Second * 60)
	Info(3333, "测试123456789")
	Error(3333, "测试12345")
	time.Sleep(time.Second * 60)
}

package internal

import (
	"bytes"
	"fmt"
	"github.com/sqc157400661/jobx/pkg/dao"
	"github.com/sqc157400661/jobx/test"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestWritingString(t *testing.T) {
	writer1 := &bytes.Buffer{}
	writer1.WriteString("111111")
	fmt.Println(writer1.String())
	writer1.WriteString("2222222")
	fmt.Println(writer1.String())
	writer1.Reset()
	go writer1.WriteString("333")
	go writer1.WriteString("4444")
	go writer1.WriteString("555")
	time.Sleep(time.Second * 10)
	fmt.Println(writer1.String())
}
func TestLogger(t *testing.T) {
	engine, err := test.GetEngine()
	dao.JFDb = engine
	require.NoError(t, err)
	logger := NewBufferLogger()
	go func() {
		for j := 0; j < 3; j++ {
			fmt.Println(j, "--")
			logger.Write(j, "hello")
		}
	}()

	go func() {
		for j := 0; j < 3; j++ {
			fmt.Println(j, "==")
			logger.Write(j, "111")
		}
	}()

	go func() {
		for i := 0; i < 3; i++ {
			fmt.Println(i, "++")
			logger.Write(i, "test")
		}
	}()

	go func() {
		for i := 0; i < 3; i++ {
			fmt.Println(i, "||")
			logger.Write(i, "test1")
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i, "$$")
			logger.Write(i, "test2")
		}
	}()

	select {}
}

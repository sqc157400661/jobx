package biz

import (
	"fmt"
	"github.com/sqc157400661/jobx/pkg/mysql"
	"github.com/sqc157400661/jobx/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateCron(t *testing.T) {
	engine, err := test.GetEngine()
	require.NoError(t, err)
	mysql.JFDb = engine
	fmt.Println(UpdateCron(1, "*/2 * * * * *"))
}

func TestDeleteCronByID(t *testing.T) {
	engine, err := test.GetEngine()
	require.NoError(t, err)
	mysql.JFDb = engine
	fmt.Println(DeleteCronByID(1))
}

func TestRebootCronByID(t *testing.T) {
	engine, err := test.GetEngine()
	require.NoError(t, err)
	mysql.JFDb = engine
	fmt.Println(RebootCronByID(2))
}

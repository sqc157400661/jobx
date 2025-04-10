package test

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/pkg/mysql"
)

func GetEngine() (*xorm.Engine, error) {
	engine, err := mysql.NewMySQLEngine(config.MySQL{
		Host:   "localhost",
		User:   "root",
		Passwd: "157400661",
		DB:     "task_center",
		Port:   3306,
	}, true, true)
	if err != nil {
		return nil, err
	}
	return engine, nil
}

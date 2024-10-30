package test

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func GetEngine() (*xorm.Engine, error) {
	engine, err := NewEngine(MySQL{
		Host:     "localhost",
		User:     "root",
		Password: "157400661",
		DB:       "pgpaas",
		Port:     3306,
	}, false)
	if err != nil {
		return nil, err
	}
	return engine, nil
}

type MySQL struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
	Port     int    `yaml:"port"`
}

func (m MySQL) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&interpolateParams=true", m.User, m.Password, m.Host, m.Port, m.DB)
}
func NewEngine(msql MySQL, verbose bool) (engine *xorm.Engine, err error) {
	engine, err = xorm.NewEngine("mysql", msql.String())
	if err != nil {
		return nil, err
	}
	err = engine.Ping()
	if err != nil {
		return nil, err
	}
	engine.ShowSQL(verbose)

	return engine, nil
}

package mysql

import (
	"github.com/go-xorm/xorm"
	"github.com/sqc157400661/jobx/config"
	"sync"
)

var setOnce sync.Once

// JFDb When the tenant appName is the same,
// only one database is allowed, so a global variable is set here
var JFDb *xorm.Engine

func DB() *xorm.Engine {
	return JFDb
}

func SetDB(conf config.MySQL) error {
	engine, err := NewMySQLEngine(conf, true, true)
	if err != nil {
		return err
	}
	setOnce.Do(func() {
		JFDb = engine
	})
	return nil
}

// DSNProvidor provides DataSourceName
type DSNProvidor interface {
	DSN() string
}

// NewMySQLEngine returns db engine of mysql
func NewMySQLEngine(mysql DSNProvidor, ping, verbose bool) (engine *xorm.Engine, err error) {
	engine, err = xorm.NewEngine("mysql", mysql.DSN())
	if err != nil {
		return
	}
	if ping {
		err = engine.Ping()
		if err != nil {
			return
		}
	}
	engine.ShowSQL(verbose)
	//engine.SetMaxIdleConns(1)
	//engine.DB().SetMaxOpenConns(1)
	return
}

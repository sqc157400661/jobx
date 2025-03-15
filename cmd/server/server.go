package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/pflag"
	"github.com/sqc157400661/helper/conf"
	"github.com/sqc157400661/helper/mysql"
	"github.com/sqc157400661/util"

	"github.com/sqc157400661/jobx/api/router"
	"github.com/sqc157400661/jobx/cmd/service"
	"github.com/sqc157400661/jobx/hack/demo"
	"github.com/sqc157400661/jobx/pkg/dao"
	"github.com/sqc157400661/jobx/pkg/providers"
)

var configFile = pflag.StringP("config", "c", "./config.yaml", "Input Config File")

func StartServer() (err error) {
	pflag.Parse()
	//初始化配置文件
	conf.Setup(*configFile)

	err = initDB()
	if err != nil {
		util.PrintFatalError(err)
	}
	//路由配置
	routersInit := router.InitRouter()
	//http服务器配置
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", conf.GetIntD("server.port", 8080)),
		Handler:        routersInit,
		ReadTimeout:    conf.GetDurationD("server.readTimeout", 30) * time.Millisecond,
		WriteTimeout:   conf.GetDurationD("server.writeTimeout", 30) * time.Millisecond,
		MaxHeaderBytes: 1 << 20,
	}
	jb, _ := service.NewJobFlow("sqc_test_compute", dao.JFDb)
	_ = jb.Register(
		&providers.DemoTasker{},
		&providers.Demo2Tasker{},
		&demo.CheckIdle{},
		&demo.MarkVwStatusInDB{},
		&demo.MarkVwPendingStatusInDB{},
		&demo.QueryCnchPendingTask{},
		&demo.QueryMetric{},
		&demo.UpdateK8sResource{},
		&demo.UpdateK8sResourceCheckLoop{},
		&demo.PreVwCheckTasker{},
	)
	jb.Start()
	//启动服务器
	err = s.ListenAndServe()
	return
}

type DatabaseConfig struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	User   string `yaml:"user"`
	Verify string `yaml:"verify"`
	DB     string `yaml:"db"`
}

func initDB() (err error) {
	var dbConf mysql.ConnectInfo
	err = conf.UnmarshalKey("mysql", &dbConf)
	if err != nil {
		return err
	}
	dao.JFDb, err = mysql.NewMySQLEngine(dbConf, true, true)
	if err != nil {
		return
	}
	return nil
}

package server

import (
	"fmt"
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/pkg/mysql"
	"net/http"
	"time"

	"github.com/sqc157400661/jobx/api/router"
)

func StartServer(config config.ServerConfig) (err error) {
	err = mysql.SetDB(config.MySQL)
	if err != nil {
		return err
	}
	//路由配置
	routersInit := router.InitRouter()
	//http服务器配置
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.ServerPort),
		Handler:        routersInit,
		ReadTimeout:    30 * time.Millisecond,
		WriteTimeout:   30 * time.Millisecond,
		MaxHeaderBytes: 1 << 20,
	}
	//启动服务器
	err = s.ListenAndServe()
	return
}

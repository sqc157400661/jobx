package main

import (
	"github.com/spf13/pflag"
	"github.com/sqc157400661/jobx/cmd/server"
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/util"
	"gopkg.in/yaml.v3"
	"os"
)

var configFile = pflag.StringP("config", "c", "./config.yaml", "Input Config File")

func main() {
	pflag.Parse()
	// 读取文件内容
	yamlFile, err := os.ReadFile(*configFile)
	if err != nil {
		util.PrintFatalError(err)
	}

	// 解析YAML
	var serverConf config.ServerConfig
	err = yaml.Unmarshal(yamlFile, &serverConf)
	if err != nil {
		util.PrintFatalError(err)
	}
	err = server.StartServer(serverConf)
	if err != nil {
		util.PrintFatalError(err)
	}
}

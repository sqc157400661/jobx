package config

import "fmt"

type ServerConfig struct {
	MySQL      MySQL  `yaml:"mysql"`
	ServerUid  string `yaml:"serverUid"`
	ServerPort int    `yaml:"serverPort"`
}

type MySQL struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	User   string `yaml:"user"`
	Passwd string `yaml:"passwd"`
	DB     string `yaml:"db"`
	Socket string `yaml:"socket"`
}

// DSN returns MySQL DataSourceName
func (m MySQL) DSN() string {
	if len(m.Socket) > 0 {
		return fmt.Sprintf("%s:%s@unix(%s)/%s?charset=utf8&interpolateParams=true", m.User, m.Passwd, m.Socket, m.DB)
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&interpolateParams=true&parseTime=true", m.User, m.Passwd, m.Host, m.Port, m.DB)
}

type JobFlowServiceConfig struct {
	Uid   string `yaml:"uid"`
	MySQL MySQL
}

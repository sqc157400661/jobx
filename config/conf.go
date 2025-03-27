package config

type MySQL struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	User   string `yaml:"user"`
	Passwd string `yaml:"passwd"`
	DB     string `yaml:"db"`
}
type JobFlowServiceConfig struct {
	Uid   string `yaml:"uid"`
	MySQL MySQL
}

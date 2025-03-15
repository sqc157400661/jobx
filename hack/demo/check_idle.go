package demo

import "fmt"

type CheckIdle struct {
	Provider
}

func (t *CheckIdle) Name() string {
	return "CheckIdle"
}

func (t *CheckIdle) Output() (ctx map[string]interface{}, res map[string]interface{}, err error) {
	return
}
func (t *CheckIdle) Run(nu int) (err error) {
	for _, vw := range t.input.VWs {
		t.logger.Info(fmt.Sprintf("Vw:%s is idle", vw))
	}
	return
}

func (t *CheckIdle) Desc() string {
	return "检查是否满足空闲标准"
}

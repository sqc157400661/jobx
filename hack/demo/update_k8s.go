package demo

import "fmt"

type UpdateK8sResource struct {
	Provider
}

func (t *UpdateK8sResource) Name() string {
	return "UpdateK8sResource"
}

func (t *UpdateK8sResource) Output() (ctx map[string]interface{}, res map[string]interface{}, err error) {
	return
}
func (t *UpdateK8sResource) Run(nu int) (err error) {
	for _, vw := range t.input.VWsQueue {
		t.logger.Info(fmt.Sprintf("Start update Vw:%s", vw))
		t.logger.Info(fmt.Sprintf("Start update Vw:%s suncces", vw))
	}
	return
}

func (t *UpdateK8sResource) Desc() string {
	return "更新对应的K8s VW资源信息"
}

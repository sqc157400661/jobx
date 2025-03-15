package demo

import "fmt"

type UpdateK8sResourceCheckLoop struct {
	Provider
}

func (t *UpdateK8sResourceCheckLoop) Name() string {
	return "CheckK8sResource"
}

func (t *UpdateK8sResourceCheckLoop) Output() (ctx map[string]interface{}, res map[string]interface{}, err error) {
	return
}

func (t *UpdateK8sResourceCheckLoop) Run(nu int) (err error) {
	for _, vw := range t.input.VWsQueue {
		t.logger.Info(fmt.Sprintf("check Vw:%s k8s resource", vw))
		t.logger.Info(fmt.Sprintf("check Vw:%s k8s resource suncces", vw))
	}
	return
}

func (t *UpdateK8sResourceCheckLoop) Desc() string {
	return "watch k8s资源是否被更新完成"
}

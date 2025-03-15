package demo

import "fmt"

type PreVwCheckTasker struct {
	Provider
}

func (t *PreVwCheckTasker) Name() string {
	return "PreVwCheck"
}

func (t *PreVwCheckTasker) Output() (ctx map[string]interface{}, res map[string]interface{}, err error) {
	ctx = map[string]interface{}{
		"vwsQueue": t.input.VWsQueue,
	}
	return
}

func (t *PreVwCheckTasker) Run(nu int) (err error) {
	VwMapStatus := map[string]string{
		"vw-21000334-test-01": "Stopped",
		"vw-21000334-test-03": "Running",
		"vw-21000334-test-06": "Running",
	}
	for _, vw := range t.input.VWs {
		t.logger.Info(fmt.Sprintf("Start Check Vw:%s", vw))
		t.logger.Info(fmt.Sprintf("Start Check success Vw:%s currentStatus:%s", vw, VwMapStatus[vw]))
		if t.input.Action == "Suspend" {
			if VwMapStatus[vw] == "Running" {
				t.logger.Info(fmt.Sprintf("Vw:%s added suspend queue", vw))
				t.input.VWsQueue = append(t.input.VWsQueue, vw)
			}
		} else {
			if VwMapStatus[vw] == "Stopped" {
				t.logger.Info(fmt.Sprintf("Vw:%s added resume queue", vw))
				t.input.VWsQueue = append(t.input.VWsQueue, vw)
			}
		}
	}

	return
}

func (t *PreVwCheckTasker) Desc() string {
	return "校验Vw的状态"
}

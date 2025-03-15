package demo

import "fmt"

type MarkVwStatusInDB struct {
	Provider
}

func (t *MarkVwStatusInDB) Name() string {
	return "MarkVwStatusInDB"
}

func (t *MarkVwStatusInDB) Output() (ctx map[string]interface{}, res map[string]interface{}, err error) {
	return
}

func (t *MarkVwStatusInDB) Run(nu int) (err error) {
	newStatus := "Running"
	if t.input.Action == "Suspend" {
		newStatus = "Stopped"
	}
	for _, vw := range t.input.VWsQueue {
		t.logger.Info(fmt.Sprintf("Start update Vw:%s new status %s in db", vw, newStatus))
		t.logger.Info(fmt.Sprintf("Start update Vw:%s  new status %s in db suncces", vw, newStatus))
	}
	return
}

func (t *MarkVwStatusInDB) Desc() string {
	return "更新vw在数据库中的状态信息"
}

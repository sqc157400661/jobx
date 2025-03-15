package demo

import "fmt"

type MarkVwPendingStatusInDB struct {
	Provider
}

func (t *MarkVwPendingStatusInDB) Name() string {
	return "LockVwStatusInDB"
}

func (t *MarkVwPendingStatusInDB) Output() (ctx map[string]interface{}, res map[string]interface{}, err error) {
	return
}

func (t *MarkVwPendingStatusInDB) Run(nu int) (err error) {
	newStatus := "ToRunning"
	if t.input.Action == "Suspend" {
		newStatus = "ToStop"
	}
	for _, vw := range t.input.VWsQueue {
		t.logger.Info(fmt.Sprintf("lock Vw:%s new status %s in db", vw, newStatus))
		t.logger.Info(fmt.Sprintf("lock Vw:%s  new status %s in db suncces", vw, newStatus))
	}
	return
}

func (t *MarkVwPendingStatusInDB) Desc() string {
	return "lock vw status in db"
}

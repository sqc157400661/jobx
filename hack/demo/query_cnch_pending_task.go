package demo

type QueryCnchPendingTask struct {
	Provider
}

func (t *QueryCnchPendingTask) Name() string {
	return "QueryPendingTask"
}

func (t *QueryCnchPendingTask) Output() (ctx map[string]interface{}, res map[string]interface{}, err error) {
	ctx = map[string]interface{}{
		"idlesVWs": []string{
			"vw-21000334-test-01",
		}}
	return
}

func (t *QueryCnchPendingTask) Run(nu int) (err error) {
	t.logger.Info("query pending task start")
	t.logger.Info("query pending task success")
	return
}

func (t *QueryCnchPendingTask) Desc() string {
	return "查询cnch的pengding-task表"
}

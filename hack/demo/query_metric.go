package demo

import "fmt"

type QueryMetric struct {
	Provider
}

func (t *QueryMetric) Name() string {
	return "QueryIdlesMetric"
}

func (t *QueryMetric) Output() (ctx map[string]interface{}, res map[string]interface{}, err error) {
	ctx = map[string]interface{}{
		"vws": []string{
			"vw-21000334-test-01",
			"vw-21000334-test-03",
			"vw-21000334-test-06",
		},
	}
	return
}

func (t *QueryMetric) Run(nu int) (err error) {
	t.logger.Info("query metric start")
	t.logger.Info(fmt.Sprintf("query PromSQL :%s", "((max(max_over_time(cnch_current_metrics_query{cluster=\"{{.Cluster}}\", pod=~\"vw.*\"}[{{.Window}}])) by (pod, cluster, namespace)) == 0) ` \n\t\t`* on (pod, namespace) group_left(workload) max(cnch:vw:metrics:workload_vw_pod:kube_pod_owner:relabel{}) by (pod, namespace)"))
	t.logger.Info("query metric success")
	return
}

func (t *QueryMetric) Desc() string {
	return "查询query、pengding-task以及backup-task的空闲指标"
}

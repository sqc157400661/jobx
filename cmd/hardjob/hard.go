package hardjob

import (
	"github.com/spf13/viper"
	"github.com/sqc157400661/jobx/api/types"
	"github.com/sqc157400661/jobx/cmd/job"
	"github.com/sqc157400661/jobx/internal/helper"
	"github.com/sqc157400661/jobx/pkg/dao"
	"github.com/sqc157400661/jobx/pkg/errors"
	"github.com/sqc157400661/jobx/pkg/options"
	"strings"
)

// NewHardJob quickly create tasks based on predefined jobs,
// which are declared and defined through the job_definition table.
// job_definition:
// job_definition.name,specify the name of the predefined job.
// job_definition.tenant,specify tenants for predefined job.
// job_definition.sort,determine the sorting rules for the job pipelines.
// job_definition.pipelines,define pipeline nodes in the format of action or name: action,
// where the action corresponds to the Name of the provider, in the format of '{"Create Runner'}`.
// job_definition.retry,number of retries.
// job_definition.input,predefined input parameters.
// job_definition.env,predefined environmental parameters.
// job_definition.condition,predefined conditions.
// determine based on the path and corresponding value of the input parameter,
// in the format of {"create.pg.pvlEnable": "true"}.
func NewHardJob(name, owner, tenant string, input interface{}, opts ...options.JobOptionFunc) (jobId int, err error) {
	inputMap, _ := helper.Struct2Map(input)
	var jobDefs []*dao.JobDefinition
	err = dao.JFDb.Where("name=? and tenant=?", name, tenant).OrderBy("sort asc").Find(&jobDefs)
	if err != nil {
		return
	}
	if len(jobDefs) == 0 {
		err = errors.NotFoundDefJob()
		return
	}
	if len(opts) > 0 {
		opts = append(opts, options.JobInput(inputMap))
	} else {
		opts = []options.JobOptionFunc{options.JobInput(inputMap)}
	}
	vipConfig := viper.New()
	defer func() { vipConfig = nil }()
	err = vipConfig.MergeConfigMap(inputMap)
	if err != nil {
		return
	}
	jober := job.NewJober(name, owner, tenant, opts...)
	for _, item := range jobDefs {
		if len(item.Pipelines) == 0 {
			continue
		}
		needAdd := true
		if len(item.Condition) > 0 {
			for key, val := range item.Condition {
				if vipConfig.GetString(key) != val {
					needAdd = false
				}
			}
		}
		if needAdd {
			addPipelines(jober, item)
		}
	}
	err = jober.Exec()
	if errors.IgnoreBIZExist(err) != nil {
		if errors.IgnoreTokenExist(err) != nil {
			if errors.IgnoreWaitTimeout(err) != nil {
				return 0, err
			}
		}
	}
	err = nil
	jobId = jober.Job.ID
	return
}

func addPipelines(jober *job.Jober, jf *dao.JobDefinition) {
	var opts []options.JobOptionFunc
	if jf.Retry > 0 {
		opts = append(opts, options.RetryNum(jf.Retry))
	}
	if len(jf.Input) > 0 {
		opts = append(opts, options.JobInput(jf.Input))
	}
	if len(jf.Env) > 0 {
		opts = append(opts, options.JobEnv(jf.Env))
	}
	for _, pipe := range jf.Pipelines {
		if pipe == "" {
			continue
		}
		var name, action = getPipelineInfo(pipe)
		jober.AddPipeline(name, action, opts...)
	}
}

func getPipelineInfo(pipe string) (name, action string) {
	pipe = strings.ReplaceAll(pipe, " ", "")
	pipeArr := strings.Split(pipe, ":")
	if len(pipeArr) == 1 {
		name = pipeArr[0]
		action = pipeArr[0]
	} else {
		name = pipeArr[0]
		action = pipeArr[1]
	}
	return
}

func GetHardJob(name, tenant string) (hardJobDefinition types.HardJobDefinition, err error) {
	var jobDefs []*dao.JobDefinition
	err = dao.JFDb.Where("name=? and tenant=?", name, tenant).OrderBy("sort asc").Find(&jobDefs)
	if err != nil {
		return
	}
	hardJobDefinition = types.HardJobDefinition{
		Name:   name,
		Tenant: tenant,
	}
	var hardJobPipelines []types.HardJobPipeline
	for _, item := range jobDefs {
		for _, pipe := range item.Pipelines {
			var pipeName, pipeAction = getPipelineInfo(pipe)
			hardJobPipelines = append(hardJobPipelines, types.HardJobPipeline{
				Name:      pipeName,
				Action:    pipeAction,
				Input:     item.Input,
				Retry:     item.Retry,
				Env:       item.Env,
				Condition: item.Condition,
			})
		}
	}
	hardJobDefinition.Pipelines = hardJobPipelines
	return
}

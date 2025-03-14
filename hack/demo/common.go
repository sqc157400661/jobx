package demo

import (
	"code.alipay.com/dbplatform/pgoperator/pkg/helper"
	"code.alipay.com/dbplatform/pgoperator/pkg/paas/common/names"
	"code.alipay.com/dbplatform/pgoperator/pkg/paas/dao"
	"code.alipay.com/dbplatform/pgoperator/pkg/paas/types/dto/instance"
	"code.alipay.com/shiqingchuang.sqc/jobflow/cmd/log"
	"code.alipay.com/shiqingchuang.sqc/jobflow/pkg/providers"
	"code.alipay.com/shiqingchuang.sqc/tools/pkg/logging"
	"fmt"
	"github.com/go-xorm/xorm"
	log2 "github.com/sqc157400661/jobx/cmd/log"
)

type Provider struct {
	input  *instance.CreateInstance
	logger log2.LoggerAdapter
}

func (m *Provider) Input(i providers.Inputer) (err error) {
	if i == nil {
		return
	}
	m.input = &instance.CreateInstance{}
	if err = helper.ConvertToStruct(i.GetInput(), m.input); err != nil {
		return
	}
	m.logger = log.NewTaskLogger(i.GetTaskID(), logging.SugarLogger)
	return
}

type InstanceInfo struct {
	Name         string
	SigmaCluster string
	Id           string
}

func GenerateInstanceNames(clusterName string, sigmaClusters []string, instanceType string) (infos []InstanceInfo) {
	if len(sigmaClusters) == 0 || instanceType == "" {
		return
	}
	num := 1
	if instanceType == names.InstanceHAStandardClassType {
		num = 2
	}
	for i := 0; i < num; i++ {
		var info InstanceInfo
		// 如果实例在同集群，则进行特殊编号
		if len(sigmaClusters) < num {
			info = InstanceInfo{
				Name:         fmt.Sprintf("%s-%d", clusterName, i),
				SigmaCluster: sigmaClusters[0],
				Id:           fmt.Sprintf("%d", i+1),
			}
		} else {
			info = InstanceInfo{
				Name:         clusterName,
				SigmaCluster: sigmaClusters[i],
				Id:           fmt.Sprintf("%d", i+1),
			}
		}
		infos = append(infos, info)
	}
	return
}

func getReplicaInstancesByClusterName(session *xorm.Session, clusterName string) (instances []dao.DatabaseInstance, err error) {
	err = session.Where("cluster_name=? and role=? and status=?", clusterName, names.MySQLReplicaRole, names.StatusRunning).Find(&instances)
	if err != nil {
		return
	}
	if len(instances) == 0 {
		err = fmt.Errorf("not found replica instance")
	}
	return
}

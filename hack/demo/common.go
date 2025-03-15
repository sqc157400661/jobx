package demo

import (
	"github.com/sqc157400661/helper/logging"
	"github.com/sqc157400661/util"

	"github.com/sqc157400661/jobx/cmd/log"
	"github.com/sqc157400661/jobx/pkg/providers"
)

type TestInput struct {
	VWs      []string
	VWsQueue []string
	Action   string
}

type Provider struct {
	input  *TestInput
	logger log.LoggerAdapter
}

func (m *Provider) Input(i providers.Inputer) (err error) {
	if i == nil {
		return
	}
	m.input = &TestInput{}
	if err = util.ConvertToStruct(i.GetInput(), m.input); err != nil {
		return
	}
	m.logger = log.NewTaskLogger(i.GetTaskID(), logging.SugarLogger)
	return
}

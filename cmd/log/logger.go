package log

import (
	"fmt"
	"go.uber.org/zap"
)

// 日志适配器
type LoggerAdapter interface {
	Info(string)
	Infof(format string, a ...interface{})
	Error(string)
	Errorf(format string, a ...interface{})
}

type JobFlowLogAdapter struct {
	TaskId int
	logger *zap.SugaredLogger
}

func NewTaskLogger(taskId int, logger ...*zap.SugaredLogger) *JobFlowLogAdapter {
	res := &JobFlowLogAdapter{
		TaskId: taskId,
	}
	if len(logger) > 0 {
		res.logger = logger[0]
		res.logger.With(zap.Int("taskId", taskId))
	}
	return res
}

func (j *JobFlowLogAdapter) Info(msg string) {
	if j.logger == nil {
		fmt.Println(msg)
	} else {
		j.logger.Info(msg)
	}
	Info(j.TaskId, msg)
}
func (j *JobFlowLogAdapter) Infof(format string, a ...interface{}) {
	if j.logger == nil {
		fmt.Printf(format, a)
	} else {
		j.logger.Infof(format, a)
	}
	Info(j.TaskId, fmt.Sprintf(format, a...))
}
func (j *JobFlowLogAdapter) Error(msg string) {
	if j.logger == nil {
		fmt.Println(msg)
	} else {
		j.logger.With(zap.Bool("jobExecErr", true)).Error(msg)
	}
	Error(j.TaskId, msg)
}
func (j *JobFlowLogAdapter) Errorf(format string, a ...interface{}) {
	if j.logger == nil {
		fmt.Printf(format, a)
	} else {
		j.logger.With(zap.Bool("jobExecErr", true)).Errorf(format, a)
	}
	Error(j.TaskId, fmt.Sprintf(format, a...))
}

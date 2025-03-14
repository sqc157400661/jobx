package demo

import (
	"errors"
)

type DemoTasker struct {
}

func (t *DemoTasker) Name() string {
	return "demo"
}

func (t *DemoTasker) Input(i Inputer) (err error) {
	return
}

func (t *DemoTasker) Output() (ctx map[string]interface{}, res map[string]interface{}, err error) {
	return
}
func (t *DemoTasker) Run(nu int) (err error) {
	return
}

func (t *DemoTasker) Desc() string {
	return "用于测试的样例任务"
}

type Demo2Tasker struct {
}

func (t *Demo2Tasker) Name() string {
	return "demo2"
}

func (t *Demo2Tasker) Input(i Inputer) (err error) {
	return
}

func (t *Demo2Tasker) Output() (ctx map[string]interface{}, res map[string]interface{}, err error) {
	res = map[string]interface{}{
		"sds": 12345,
	}
	ctx = map[string]interface{}{
		"uid": "12334534543",
	}
	return
}
func (t *Demo2Tasker) Run(nu int) (err error) {
	return errors.New("run err for demo2")
}

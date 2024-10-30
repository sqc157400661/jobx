package providers

import (
	"time"
)

type DelayTasker struct {
	timer time.Duration
}

func (t *DelayTasker) Name() string {
	return "delay"
}

func (t *DelayTasker) Input(i Inputer) (err error) {
	input := i.GetInput()
	t.timer = 10 * time.Second
	if v, has := input["time"]; has {
		timer, ok := v.(float64)
		if ok {
			t.timer = time.Duration(timer)
		}
	}
	return
}

func (t *DelayTasker) Output() (ctx map[string]interface{}, res map[string]interface{}, err error) {
	res = map[string]interface{}{"wait_time": t.timer}
	return
}

func (t *DelayTasker) Run(nu int) (err error) {
	if t.timer > 0 {
		time.Sleep(t.timer)
	}
	return
}

func (t *DelayTasker) Desc() string {
	return "用于控制任务之间，需要等待的时间时长，返回结果里的wait_time，代表需要等待的时间"
}

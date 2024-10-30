package options

import (
	"fmt"
	"github.com/sqc157400661/jobx/config"
	"time"
)

const (
	RetryGapKey           = "retry_gap_second"
	DefaultRetryGapSecond = 10
)

type JobOptionFunc func(o *jobOptions)

// options is an application options.
type jobOptions struct {
	Desc        string
	Env         map[string]interface{}
	Pause       int8
	Retry       int
	Input       map[string]interface{}
	BizId       string
	Tokens      []string
	PreLockUid  string
	Sync        bool
	SyncTimeOut int
}

var DefaultJobOptions = jobOptions{
	Desc:  "普通任务",
	Pause: 1,
	Retry: 3,
	Env: map[string]interface{}{
		RetryGapKey: DefaultRetryGapSecond,
	},
}

// 可行参数，指定任务的描述信息
func JobDesc(desc string) JobOptionFunc {
	return func(o *jobOptions) { o.Desc = desc }
}

// 可选参数，指定入参数据
func JobInput(input map[string]interface{}) JobOptionFunc {
	return func(o *jobOptions) {
		if len(input) > 0 {
			o.Input = input
		}
	}
}

// 可选参数，指定环境参数数据
func JobEnv(env map[string]interface{}) JobOptionFunc {
	return func(o *jobOptions) {
		if len(env) > 0 {
			o.Env = env
		}
	}
}

// 可选参数，指定业务唯一编码，也是任务的幂等编码
func BizId(id string) JobOptionFunc {
	return func(o *jobOptions) { o.BizId = id }
}

// 可选参数，指定任务执行者
func PreLock(uid string) JobOptionFunc {
	return func(o *jobOptions) {
		if uid != "" {
			o.PreLockUid = fmt.Sprintf("%s%s", config.PreLockPrefix, uid)
		}
	}
}

// 可选参数，指定令牌，同类任务令牌是互斥的，只有当一个任务完成或则废弃后令牌释放后，下一个任务才可以执行
func AddTokens(tokens []string) JobOptionFunc {
	return func(o *jobOptions) {
		o.Tokens = tokens
	}
}

// 可选参数，指定该任务或者流水线节点是否可以执行暂停的操作
func Pause(t bool) JobOptionFunc {
	return func(o *jobOptions) {
		if t {
			o.Pause = 1
		} else {
			o.Pause = 0
		}
	}
}

// 可选参数，指定该任务或者流水线节点不可暂停的操作
func NoPause() JobOptionFunc {
	return func(o *jobOptions) { o.Pause = 0 }
}

// 可选参数，指定该任务或者流水线节点重试的最大次数
func RetryNum(retry int) JobOptionFunc {
	return func(o *jobOptions) { o.Retry = retry }
}

// 可选参数，指定该任务或者流水线重试的时间间隔梯度
func RetryGapSecond(t int) JobOptionFunc {
	return func(o *jobOptions) {
		if o.Env == nil {
			o.Env = map[string]interface{}{
				RetryGapKey: t,
			}
		} else {
			o.Env[RetryGapKey] = t
		}
	}
}

// 可选参数，指定该是否等待任务的执行结果
func Sync(sync bool, timeout int) JobOptionFunc {
	if sync && timeout == 0 {
		timeout = 5
	}
	return func(o *jobOptions) {
		o.Sync = sync
		o.SyncTimeOut = timeout
	}
}

func GetRetryGapSecond(env map[string]interface{}) (t time.Duration) {
	t = DefaultRetryGapSecond * time.Second
	if env == nil {
		return
	}
	if v, has := env[RetryGapKey]; has {
		gap, ok := v.(int)
		if ok {
			t = time.Duration(gap) * time.Second
		}
	}
	return
}

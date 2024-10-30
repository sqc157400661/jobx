package options

import (
	"time"
)

// Option is an application option.
type OptionFunc func(o *Options)

// options is an application options.
type Options struct {
	// desc for JobFlow
	Desc string
	//Timeout      int // 单位秒，执行超时时间
	// LoopInterval control the time interval of cyclic data fetching
	LoopInterval time.Duration
	// PoolLen control the number of tasks executed at the same time
	PoolLen int
}

var DefaultOption = Options{
	Desc:         "JobFlow run for flow job",
	PoolLen:      2,
	LoopInterval: 5 * time.Second,
}

// Name with JobFlow desc.
func Desc(desc string) OptionFunc {
	return func(o *Options) { o.Desc = desc }
}

// LoopInterval with JobFlow LoopInterval.
func LoopInterval(t time.Duration) OptionFunc {
	return func(o *Options) { o.LoopInterval = t }
}

// PoolLen with JobFlow PoolLen.
func PoolLen(l int) OptionFunc {
	return func(o *Options) { o.PoolLen = l }
}

package cronopt

import "github.com/sqc157400661/jobx/config"

// CronOptionFunc is an application option.
type CronOptionFunc func(o *CronOptions)

// CronOptions is an application options.
type CronOptions struct {
	Tenant         string
	AppName        string
	CurrencyPolicy string
	SecondEnable   bool
}

var DefaultOption = CronOptions{
	CurrencyPolicy: config.AllowCronCurrencyPolicy,
}

// AppName with JobFlow AppName.
func AppName(l string) CronOptionFunc {
	return func(o *CronOptions) { o.AppName = l }
}

// Tenant with JobFlow Tenant.
func Tenant(l string) CronOptionFunc {
	return func(o *CronOptions) { o.Tenant = l }
}

func CurrencyAllow() CronOptionFunc {
	return func(o *CronOptions) { o.CurrencyPolicy = config.AllowCronCurrencyPolicy }
}

func CurrencyForbid() CronOptionFunc {
	return func(o *CronOptions) { o.CurrencyPolicy = config.ForbidCronCurrencyPolicy }
}

func CurrencyReplace() CronOptionFunc {
	return func(o *CronOptions) { o.CurrencyPolicy = config.ReplaceCronCurrencyPolicy }
}

func Second() CronOptionFunc {
	return func(o *CronOptions) { o.SecondEnable = true }
}

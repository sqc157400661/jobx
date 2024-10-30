package errors

const (
	// Code for error

	NoError     = 0
	ServerError = 499
	ParamError  = 400
	// for Tracker
	// for Worker
	// for DB
)

var ErrorReasonMap = map[int]string{
	BIZExist:       BIZExistReason,
	TokenExist:     TokenExistReason,
	NoJobProvider:  NoJobProviderReason,
	WaitJobTimeout: WaitJobTimeoutReason,
	UnknownCode:    UnknownReason,
}

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = 5001
	// UnknownReason is unknown reason for error info.
	UnknownReason = "未知错误"

	BIZExist         = 1001
	BIZExistReason   = "唯一编码冲突"
	TokenExist       = 1002
	TokenExistReason = "令牌冲突，相关任务可能已经在执行中"

	NoJobProvider       = 2001
	NoJobProviderReason = "没有找到任务执行者"

	WaitJobTimeout       = 3001
	WaitJobTimeoutReason = "任务在等待时间内未完成"

	NotFoundJobDefs       = 4001
	NotFoundJobDefsReason = "暂不支持该操作"
)

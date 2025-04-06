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
	BIZExist:                   BIZExistReason,
	TokenExist:                 TokenExistReason,
	NoJobProvider:              NoJobProviderReason,
	WaitJobTimeout:             WaitJobTimeoutReason,
	UnknownCode:                UnknownReason,
	LocalQueueFullCode:         LocalQueueFullReason,
	ErrEmptyQueueCode:          ErrEmptyQueueReason,
	ErrTaskNotFoundCode:        ErrTaskNotFoundReason,
	NotFoundJobDefs:            NotFoundJobDefsReason,
	NoTaskFoundCode:            NoTaskFoundReason,
	NoRunnablePipelineTaskCode: NoRunnablePipelineTaskReason,
	InvalidParameterCode:       InvalidParameterReason,
	JobDefExistCode:            JobDefExistReason,
	CronJobExistCode:           CronJobExistReason,
}

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = 100
	// UnknownReason is unknown reason for error info.
	UnknownReason          = "Unknown error"
	InvalidParameterCode   = 101
	InvalidParameterReason = "Invalid parameter"

	BIZExist           = 1001
	BIZExistReason     = "biz conflict,biz uid already exists"
	TokenExist         = 1002
	TokenExistReason   = "Token conflict, related tasks may already be in progress"
	JobDefExistCode    = 1003
	JobDefExistReason  = "job def exist"
	CronJobExistCode   = 1004
	CronJobExistReason = "cronjob exist"

	NoJobProvider                = 2001
	NoJobProviderReason          = "Task provider not found"
	NoTaskFoundCode              = 2002
	NoTaskFoundReason            = "no task found"
	NoRunnablePipelineTaskCode   = 2004
	NoRunnablePipelineTaskReason = "no runnable pipeline task"

	LocalQueueFullCode    = 2101
	LocalQueueFullReason  = "task queue is at full capacity"
	ErrEmptyQueueCode     = 2102
	ErrEmptyQueueReason   = "no tasks in queue"
	ErrTaskNotFoundCode   = 2103
	ErrTaskNotFoundReason = "task not found in queue"

	WaitJobTimeout       = 3001
	WaitJobTimeoutReason = "The task was not completed within the waiting time"

	NotFoundJobDefs       = 4001
	NotFoundJobDefsReason = "This operation is not supported"
)

var (
	ErrInvalidParameter = New(InvalidParameterCode, "")

	// for Job
	ErrNoProvider     = New(NoJobProvider, "")
	ErrNotFoundDefJob = New(NotFoundJobDefs, "")
	ErrWaitJobTimeout = New(WaitJobTimeout, "")

	// for local queue
	ErrQueueFullError    = New(LocalQueueFullCode, "")
	ErrEmptyQueueError   = New(ErrEmptyQueueCode, "")
	ErrTaskNotFoundError = New(ErrTaskNotFoundCode, "")

	ErrJobDefExist  = New(JobDefExistCode, "")
	ErrCronJobExist = New(CronJobExistCode, "")
)

package config

const (
	MaxJobLevel = 2

	// status  used externally
	StatusPending   = "pending"
	StatusPause     = "pause"
	StatusSkip      = "skip"
	StatusFail      = "fail"
	StatusDiscarded = "discarded"
	StatusSuccess   = "success"

	// phase   used internally
	PhaseInit       = "init"
	PhaseReady      = "ready"
	PhaseRunning    = "running"
	PhaseTerminated = "terminated"

	RunnableYes = 1
	RunnableNo  = 0

	// cron
	CronStatusValid     = "valid"     // 有效的，待运行
	CronStatusRunning   = "running"   // 运行中的
	CronStatusRebooting = "rebooting" // 重启中的
	CronStatusUpdating  = "updating"  // 更新中的
	CronStatusInvalid   = "invalid"   // 无效的
	CronStatusDeleted   = "deleted"   // 已删除的
	CronExecFuncType    = "func"      // func任务类型
	CronExecJobType     = "job.run"   // job任务类型

	// Pre-Locking
	PreLockPrefix = "pre-"
)

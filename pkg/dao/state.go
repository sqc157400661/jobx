package dao

import (
	"github.com/sqc157400661/jobx/config"
)

type State struct {
	Phase  string `gorm:"column:phase" json:"phase" xorm:"phase"`    // 状态控制
	Status string `gorm:"column:status" json:"status" xorm:"status"` // 展示控制和手工控制
	Reason string `gorm:"column:reason" json:"reason" xorm:"reason"`
}

// IsReady 是否可以执行
func (s State) IsReady() bool {
	if s.Phase == config.PhaseReady && s.Status == config.StatusPending {
		return true
	}
	return false
}

// IsReady 是否可以执行
func (s State) IsRunning() bool {
	if s.Phase == config.PhaseRunning {
		return true
	}
	return false
}

// IsFail 是否失败
func (s State) IsFailed() bool {
	if s.Status == config.StatusFail {
		return true
	}
	return false
}

// IsFinished 是否已经完成
func (s State) IsFinished() bool {
	if s.IsSuccess() || s.IsSkip() {
		return true
	}
	return false
}

// IsSuccess 是否已经成功执行
func (s State) IsSuccess() bool {
	if s.Status == config.StatusSuccess {
		return true
	}
	return false
}

// IsDiscarded 是否废弃
func (s State) IsDiscarded() bool {
	if s.Status == config.StatusDiscarded {
		return true
	}
	return false
}

// IsPausing 是否暂停中
func (s State) IsPausing() bool {
	if s.Status == config.StatusPause {
		return true
	}
	return false
}

// IsSkip 是否要执行跳过
func (s State) IsSkip() bool {
	if s.Status == config.StatusSkip {
		return true
	}
	return false
}

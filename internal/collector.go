package internal

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/pkg/dao"
)

type Collector interface {
	loadCheckUndoJobs(uid string) (jobs []dao.Job, err error)
	CollectByUID(uid string) (jobs []*dao.Job, err error)
	Release(uid string) (err error)
}

type DefaultCollector struct {
	len int
}

func NewDefaultCollector(len int) (collector *DefaultCollector) {
	if len <= 0 {
		len = 1
	}
	return &DefaultCollector{
		len: len,
	}
}

// 启动初始化的时候检查未完成的job队列
func (c *DefaultCollector) loadCheckUndoJobs(uid string) (jobs []dao.Job, err error) {
	err = dao.JFDb.Where("locker=? and phase !=? and parent_id=0", uid, config.PhaseTerminated).Find(&jobs)
	return
}

func (c *DefaultCollector) CollectByUID(uid string) (jobs []*dao.Job, err error) {
	var num int64
	num, err = c.lock(uid)
	if err != nil || num == 0 {
		return
	}
	err = dao.JFDb.Where("locker=? and phase =? and parent_id=0", uid, config.PhaseReady).Find(&jobs)
	return
}

func (c *DefaultCollector) Release(uid string) (err error) {
	var jobs []dao.Job
	err = dao.JFDb.In("phase", []string{config.PhaseReady, config.PhaseRunning}).Where("locker=?", uid).Find(&jobs)
	if err != nil {
		return
	}
	// 依次更新状态并解除锁定
	for _, v := range jobs {
		err = c.unLock(uid, v.ID)
		if err != nil {
			err = errors.Wrapf(err, "uid:%s unlock err id:%d", uid, v.ID)
			return
		}
	}
	return
}

// AddJobLocker 对job添加locker
func (c *DefaultCollector) lock(uid string) (lockedNum int64, err error) {
	var res sql.Result
	// CAS
	uidPreLock := fmt.Sprintf("%s%s", config.PreLockPrefix, uid)
	res, err = dao.JFDb.Exec(`update job set locker=? where  parent_id=0 and (locker='' or locker=?) and phase =? order by id asc limit ?`, uid, uidPreLock, config.PhaseReady, c.len)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (c *DefaultCollector) unLock(uid string, jobId int) (err error) {
	_, err = dao.JFDb.Exec(`update job set locker='',phase =?  where locker=? and id =?`, config.PhaseReady, uid, jobId)
	return
}

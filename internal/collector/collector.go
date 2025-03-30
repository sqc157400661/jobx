package collector

import (
	"database/sql"
	"github.com/go-xorm/xorm"
	"github.com/pkg/errors"
	"github.com/sqc157400661/jobx/config"
	"github.com/sqc157400661/jobx/internal/names"
	"github.com/sqc157400661/jobx/pkg/dao"
)

const (
	checkUndoJobsSqlTmpl = "locker=? and phase !=? and parent_id=0"
	getStealJobSqlTmpl   = "locker=? and phase =? and parent_id=0"
	stealJobsSqlTmpl     = `update job set locker=? where  parent_id=0 and (locker='' or locker=?) and phase =? order by id asc limit ?`
	releaseJobSqlTmpl    = `update job set locker='',phase =?  where locker=? and id =?`
)

type Collector interface {
	StealJob() (jobs []*dao.Job, err error)
	ReleaseJob() (err error)
}

type DefaultCollector struct {
	stealLen  int
	serverUid string
	engine    *xorm.Engine
}

func NewDefaultCollector(engine *xorm.Engine, serverUid string, len int) (collector *DefaultCollector) {
	if len <= 0 {
		len = 1
	}
	return &DefaultCollector{
		serverUid: serverUid,
		stealLen:  len,
		engine:    engine,
	}
}

// 启动初始化的时候检查未完成的job队列
func (c *DefaultCollector) loadCheckUndoJobs(uid string) (jobs []dao.Job, err error) {
	err = c.engine.Where(checkUndoJobsSqlTmpl, uid, config.PhaseTerminated).Find(&jobs)
	return
}

func (c *DefaultCollector) StealJob() (jobs []*dao.Job, err error) {
	var num int64
	num, err = c.steal()
	if err != nil || num == 0 {
		return
	}
	err = c.engine.Where(getStealJobSqlTmpl, c.serverUid, config.PhaseReady).Find(&jobs)
	return
}

func (c *DefaultCollector) ReleaseJob() (err error) {
	var jobs []dao.Job
	err = c.engine.In("phase", []string{config.PhaseReady, config.PhaseRunning}).Where("locker=?", c.serverUid).Find(&jobs)
	if err != nil {
		return
	}
	// 依次更新状态并解除锁定
	for _, v := range jobs {
		err = c.releaseByID(v.ID)
		if err != nil {
			err = errors.Wrapf(err, "uid:%s unlock err id:%d", c.serverUid, v.ID)
			return
		}
	}
	return
}

// AddJobLocker 对job添加locker
func (c *DefaultCollector) steal() (lockedNum int64, err error) {
	var res sql.Result
	// CAS
	res, err = c.engine.Exec(stealJobsSqlTmpl, c.serverUid, names.PreLockKey(c.serverUid), config.PhaseReady, c.stealLen)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (c *DefaultCollector) releaseByID(jobId int) (err error) {
	_, err = c.engine.Exec(releaseJobSqlTmpl, config.PhaseReady, c.serverUid, jobId)
	return
}

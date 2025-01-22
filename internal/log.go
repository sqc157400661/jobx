package internal

import (
	"fmt"
	"github.com/sqc157400661/jobx/pkg/dao"
	"github.com/sqc157400661/jobx/pkg/log"
	"k8s.io/klog/v2"
	"strings"
	"sync"
	"time"
)

type loggerItem struct {
	eventID int
	msg     string
}

type Logger struct {
	maxLen      int // 最大的日志长度
	maxSetSize  int // 处理集合的最大长度,即同时可记录的任务日志数量,超过刷盘
	set         *log.LoggerSet
	flushLoop   *time.Ticker // 刷盘频率
	rebuildLoop *time.Ticker // 刷盘频率
	bufChan     chan *loggerItem
	isLock      bool
	sync.RWMutex
}

func NewLogger() *Logger {
	logger := &Logger{
		maxLen:      2048,
		maxSetSize:  20,
		flushLoop:   time.NewTicker(time.Second * 10),
		rebuildLoop: time.NewTicker(time.Hour * 12),
		bufChan:     make(chan *loggerItem, 1000),
		set:         log.NewLoggerSet(),
	}
	go func() {
		for item := range logger.bufChan {
			if logger.set.Size() >= logger.maxSetSize {
				err := logger.Flush()
				if err != nil {
					klog.Errorf("flush Err:%s", err.Error())
				}
			}
			fmt.Printf("id:%d  msg:%s \n", item.eventID, item.msg)
			if logger.isLocked() {
				time.Sleep(time.Millisecond * 300)
				logger.Write(item.eventID, item.msg)
			} else {
				logger.set.AddOrGet(item.eventID).WriteString(item.msg)
			}

		}
	}()
	go func() {
		for {
			select {
			case <-logger.flushLoop.C:
				// 定时任务处理逻辑
				err := logger.Flush()
				if err != nil {
					klog.Errorf("flush Err:%s", err.Error())
				}
			}
		}
	}()
	go logger.rebuild()
	return logger
}

func (l *Logger) Write(eventID int, msg string) {
	if len(msg) > l.maxLen {
		msg = SubStrDecodeRuneInString(msg, l.maxLen) + "..."
	}
	l.bufChan <- &loggerItem{eventID: eventID, msg: msg}
}

func (l *Logger) isLocked() bool {
	l.RLock()
	defer l.RUnlock()
	return l.isLock
}

func (l *Logger) lock() {
	l.Lock()
	l.isLock = true
	l.Unlock()
}

func (l *Logger) unlock() {
	l.Lock()
	l.isLock = false
	l.Unlock()
}

func (l *Logger) rebuild() {
	for {
		select {
		case <-l.rebuildLoop.C:
			l.lock()
			newSet := log.NewLoggerSet()
			currentSet := l.set
			var data string
			for k, v := range currentSet.Items() {
				data = v.String()
				if len(data) == 0 {
					continue
				}
				newSet.ReAdd(k, v)
			}
			currentSet = nil
			l.set = newSet
			l.unlock()
		}
	}
}

func (l *Logger) Flush() (err error) {
	if l.isLocked() {
		return
	}
	currentSet := l.set
	// 对currentSet进行刷盘清理
	var data string
	newLogs := make([]*dao.JobLogs, 0)
	var jobLog dao.JobLogs
	for k, v := range currentSet.Items() {
		data = v.String()
		v.Reset()
		// 刷盘的时候发现set没有数据,说明任务已经完成，则进行clear todo 优化
		if len(data) == 0 {
			continue
		}
		// 拼接批量插入的数据
		// 查询数据库中是否存在有记录，有记录进行append
		jobLog, err = dao.GetLogByEventID(k)
		if err != nil {
			klog.Errorf("jobLog get by eventID Err:%s,eventID:%d", err.Error(), k)
			continue
		}
		if jobLog.ID > 0 {
			jobLog.Result = strings.Join([]string{jobLog.Result, data}, "")
			err = jobLog.Update()
			if err != nil {
				klog.Errorf("jobLog append update Err:%s,eventID:%d，jobLog:%+v", err.Error(), k, jobLog)
			}
			continue
		}
		// 如果没有记录则进行插入操作，组装批量插入的数据
		newLogs = append(newLogs, &dao.JobLogs{EventID: k, Result: data})
	}
	// 如果有待批量插入的数据，则进行插入操作
	if len(newLogs) > 0 {
		err = dao.BatchAddLogs(newLogs)
		if err != nil {
			klog.Errorf("BatchAddLogs Err:%s,logs:%+v", err.Error(), newLogs)
		}
	}
	return
}

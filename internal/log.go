package internal

import (
	"github.com/sqc157400661/jobx/internal/helper"
	"github.com/sqc157400661/jobx/pkg/log"
	"github.com/sqc157400661/jobx/pkg/model"
	"k8s.io/klog/v2"
	"strings"
	"sync"
	"time"
)

// loggerItem represents a single log entry in buffer
type loggerItem struct {
	eventID int
	msg     string
}

// BufferLogger handles buffered logging operations with periodic flushing
type BufferLogger struct {
	// Maximum length per log message
	maxLen int
	// Maximum buffer size before forced flush
	maxSetSize int
	// Interval for periodic flushing
	flushInterval time.Duration
	// Buffered channel for log entries
	bufChan chan *loggerItem
	// Buffered log storage
	set *log.LoggerSet
	// Ticker for periodic flushing
	flushTicker *time.Ticker
	// Mutex for concurrent access
	sync.RWMutex
}

// NewBufferLogger creates a new Logger instance with default values
func NewBufferLogger() *BufferLogger {
	logger := &BufferLogger{
		maxLen:      2048,
		maxSetSize:  50,
		flushTicker: time.NewTicker(time.Second * 10),
		bufChan:     make(chan *loggerItem, 2000),
		set:         log.NewLoggerSet(),
	}
	go logger.processBuffer()
	go logger.flushLoop()
	return logger
}

// processBuffer handles incoming log entries from buffer channel
func (l *BufferLogger) processBuffer() {
	for item := range l.bufChan {
		if l.set.Size() >= l.maxSetSize {
			err := l.Flush()
			if err != nil {
				klog.Errorf("flush Err:%s", err.Error())
			}
		}
		l.set.AddOrGet(item.eventID).WriteString(item.msg)
	}
}

// flushLoop handles periodic flushing
func (l *BufferLogger) flushLoop() {
	for {
		select {
		case <-l.flushTicker.C:
			// 定时任务处理逻辑
			err := l.Flush()
			if err != nil {
				klog.Errorf("flush Err:%s", err.Error())
			}
		}
	}
}

// Write adds a new log entry to the buffer
func (l *BufferLogger) Write(eventID int, msg string) {
	if len(msg) > l.maxLen {
		msg = helper.SubStrDecodeRuneInString(msg, l.maxLen) + "..."
	}
	l.bufChan <- &loggerItem{eventID: eventID, msg: msg}
}

func (l *BufferLogger) Flush() (err error) {
	currentSet := l.set
	// 对currentSet进行刷盘清理
	var data string
	newLogs := make([]*model.JobLogs, 0)
	var jobLog model.JobLogs
	for id, buf := range currentSet.Items() {
		data = buf.String()
		buf.Reset()
		// 刷盘的时候发现set没有数据,说明任务已经完成，则进行clear todo 优化
		if len(data) == 0 {
			continue
		}
		// 拼接批量插入的数据
		// 查询数据库中是否存在有记录，有记录进行append
		jobLog, err = model.GetLogByEventID(id)
		if err != nil {
			klog.Errorf("jobLog get by eventID Err:%s,eventID:%d", err.Error(), id)
			continue
		}
		if jobLog.ID > 0 {
			jobLog.Result = strings.Join([]string{jobLog.Result, data}, "")
			err = jobLog.Update()
			if err != nil {
				klog.Errorf("jobLog append update Err:%s,eventID:%d，jobLog:%+v", err.Error(), id, jobLog)
			}
			continue
		}
		// 如果没有记录则进行插入操作，组装批量插入的数据
		newLogs = append(newLogs, &model.JobLogs{EventID: id, Result: data})
	}
	// 如果有待批量插入的数据，则进行插入操作
	if len(newLogs) > 0 {
		err = model.BatchAddLogs(newLogs)
		if err != nil {
			klog.Errorf("BatchAddLogs Err:%s,logs:%+v", err.Error(), newLogs)
		}
	}
	currentSet.CleanupEmpty()
	return
}

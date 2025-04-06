package log

import (
	"fmt"
	"github.com/sqc157400661/jobx/internal/log"
	"strings"
	"time"
)

var Logger *log.BufferLogger

func init() {
	Logger = log.NewBufferLogger()
}

func Info(rootId interface{}, message string) {
	var eventID int
	switch rootId.(type) {
	case int:
		eventID = rootId.(int)
	case int64:
		eventID = int(rootId.(int64))
	case float64:
		eventID = int(rootId.(float64))
	default:
		return
	}
	format := fmt.Sprintf("[INFO] [%s] [%d]", time.Now().Format("2006-01-02 15:04:05"), eventID)
	message = strings.Join([]string{format, message, "\n"}, " ")
	Logger.Write(eventID, message)
}

func Error(rootId interface{}, message string) {
	var eventID int
	switch rootId.(type) {
	case int:
		eventID = rootId.(int)
	case int64:
		eventID = int(rootId.(int64))
	case float64:
		eventID = int(rootId.(float64))
	default:
		return
	}
	format := fmt.Sprintf("[ERROR] [%s] [%d]", time.Now().Format("2006-01-02 15:04:05"), eventID)
	message = strings.Join([]string{format, message, "\n"}, " ")
	Logger.Write(eventID, message)
}

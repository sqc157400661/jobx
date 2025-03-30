// Package queue implements an enhanced concurrent task queue with time tracking and ID-based removal
package queue

import (
	"container/list"
	"errors"
	"github.com/sqc157400661/jobx/pkg/dao"
	"sync"
	"time"
)

var (
	ErrQueueFull    = errors.New("task queue is at full capacity")
	ErrEmptyQueue   = errors.New("no tasks in queue")
	ErrTaskNotFound = errors.New("task not found in queue")
)

// JobTask with enqueue timestamp
type JobTask struct {
	dao.Job
	EnqueueTime time.Time // Time when task was enqueued 任务入队时间
}

// TaskQueue with enhanced functionality
// 增强功能的任务队列
type TaskQueue struct {
	maxSize      int                   // Maximum allowed pending tasks
	pending      *list.List            // List for pending tasks
	processing   map[int]*JobTask      // Map of in-progress tasks
	pendingMap   map[int]*list.Element // Map for O(1) pending task access
	mu           sync.RWMutex          // Main mutex for pending operations
	processingMu sync.RWMutex          // Mutex for processing map
}

// NewTaskQueue creates a new enhanced task queue
// 创建新的增强版任务队列
func NewTaskQueue(maxSize int) *TaskQueue {
	return &TaskQueue{
		maxSize:    maxSize,
		pending:    list.New(),
		processing: make(map[int]*JobTask),
		pendingMap: make(map[int]*list.Element),
	}
}

// AddToFront adds a task to the front of the queue with timestamp
// 添加任务到队列头部并记录时间戳
func (q *TaskQueue) AddToFront(task *JobTask) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.pending.Len() >= q.maxSize {
		return ErrQueueFull
	}

	task.EnqueueTime = time.Now()
	element := q.pending.PushFront(task)
	q.pendingMap[task.ID] = element
	return nil
}

// AddToBack adds a task to the end of the queue with timestamp
// 添加任务到队列尾部并记录时间戳
func (q *TaskQueue) AddToBack(task *JobTask) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.pending.Len() >= q.maxSize {
		return ErrQueueFull
	}

	task.EnqueueTime = time.Now()
	element := q.pending.PushBack(task)
	q.pendingMap[task.ID] = element
	return nil
}

// Dequeue moves a task from pending to processing state
// 从待处理队列取出任务并标记为执行中状态
func (q *TaskQueue) Dequeue() (*JobTask, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.pending.Len() == 0 {
		return nil, ErrEmptyQueue
	}

	element := q.pending.Front()
	task := element.Value.(*JobTask)
	q.pending.Remove(element)
	delete(q.pendingMap, task.ID)

	q.processingMu.Lock()
	defer q.processingMu.Unlock()
	q.processing[task.ID] = task

	return task, nil
}

// GetLongPendingTasks returns tasks pending longer than specified duration
// 获取等待时间超过指定时长的任务
func (q *TaskQueue) GetLongPendingTasks(duration time.Duration) []*JobTask {
	q.mu.RLock()
	defer q.mu.RUnlock()

	var result []*JobTask
	threshold := time.Now().Add(-duration)

	for e := q.pending.Front(); e != nil; e = e.Next() {
		task := e.Value.(*JobTask)
		if task.EnqueueTime.Before(threshold) {
			result = append(result, task)
		}
	}

	return result
}

// RemoveTaskByID removes a task from pending queue by ID
// 根据ID从待处理队列中移除任务
func (q *TaskQueue) RemoveTaskByID(taskID int) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	element, exists := q.pendingMap[taskID]
	if !exists {
		return ErrTaskNotFound
	}

	q.pending.Remove(element)
	delete(q.pendingMap, taskID)
	return nil
}

// CompleteTask removes a task from processing state
// 完成任务并从执行中状态移除
func (q *TaskQueue) CompleteTask(taskID int) {
	q.processingMu.Lock()
	defer q.processingMu.Unlock()
	delete(q.processing, taskID)
}

// PendingCount returns current number of pending tasks
// 获取待处理任务数量
func (q *TaskQueue) PendingCount() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.pending.Len()
}

// ProcessingCount returns current number of in-progress tasks
// 获取执行中任务数量
func (q *TaskQueue) ProcessingCount() int {
	q.processingMu.RLock()
	defer q.processingMu.RUnlock()
	return len(q.processing)
}

// GetTaskState returns the current state of a task
// 获取任务当前状态
func (q *TaskQueue) GetTaskState(taskID int) (string, bool) {
	q.processingMu.RLock()
	defer q.processingMu.RUnlock()

	_, inProcessing := q.processing[taskID]
	if inProcessing {
		return "processing", true
	}

	q.mu.RLock()
	defer q.mu.RUnlock()
	for e := q.pending.Front(); e != nil; e = e.Next() {
		if t := e.Value.(*JobTask); t.ID == taskID {
			return "pending", true
		}
	}

	return "", false
}

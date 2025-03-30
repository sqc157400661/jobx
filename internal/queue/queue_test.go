package queue

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskQueue_EnqueueWithTimestamp(t *testing.T) {
	q := NewTaskQueue(10)
	task := &PipelineTask{ID: 1}

	err := q.AddToBack(task)
	assert.NoError(t, err)

	// Verify timestamp was set
	assert.False(t, task.EnqueueTime.IsZero())
	assert.WithinDuration(t, time.Now(), task.EnqueueTime, 100*time.Millisecond)
}

func TestGetLongPendingTasks(t *testing.T) {
	q := NewTaskQueue(10)

	// Add tasks with different timestamps
	oldTask := &PipelineTask{ID: 1}
	q.AddToBack(oldTask)

	// Force old timestamp
	oldTime := time.Now().Add(-20 * time.Second)
	oldTask.EnqueueTime = oldTime

	newTask := &PipelineTask{ID: 2}
	q.AddToBack(newTask)

	// Test getting tasks older than 10s
	longPending := q.GetLongPendingTasks(10 * time.Second)
	assert.Len(t, longPending, 1)
	assert.Equal(t, 1, longPending[0].ID)

	// Test getting tasks older than 30s
	longPending = q.GetLongPendingTasks(30 * time.Second)
	assert.Empty(t, longPending)

	// Test getting tasks older than 0s (all pending)
	longPending = q.GetLongPendingTasks(0)
	assert.Len(t, longPending, 2)
}

func TestRemoveTaskByID(t *testing.T) {
	q := NewTaskQueue(10)

	// Add test tasks
	task1 := &PipelineTask{ID: 1}
	task2 := &PipelineTask{ID: 2}
	q.AddToBack(task1)
	q.AddToBack(task2)

	// Verify initial state
	assert.Equal(t, 2, q.PendingCount())

	// Remove existing task
	err := q.RemoveTaskByID(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, q.PendingCount())

	// Try to remove non-existent task
	err = q.RemoveTaskByID(999)
	assert.Equal(t, ErrTaskNotFound, err)

	// Verify remaining task
	remaining, err := q.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 2, remaining.ID)
}

func TestConcurrentTimeBasedOperations(t *testing.T) {
	q := NewTaskQueue(100)
	var wg sync.WaitGroup

	// Add tasks concurrently
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			task := &PipelineTask{ID: id}
			if id%2 == 0 {
				q.AddToFront(task)
			} else {
				q.AddToBack(task)
			}
		}(i)
	}

	// Let some time pass
	time.Sleep(100 * time.Millisecond)

	// Concurrently query long pending tasks
	wg.Add(1)
	go func() {
		defer wg.Done()
		longPending := q.GetLongPendingTasks(50 * time.Millisecond)
		assert.True(t, len(longPending) > 0)
	}()

	// Concurrently remove tasks
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = q.RemoveTaskByID(10) // Try to remove one
	}()

	wg.Wait()
	assert.Equal(t, 49, q.PendingCount()) // We removed 1 task
}

func TestTaskQueue_EdgeCases(t *testing.T) {
	q := NewTaskQueue(2)

	// Test empty queue
	longPending := q.GetLongPendingTasks(10 * time.Second)
	assert.Empty(t, longPending)

	err := q.RemoveTaskByID(999)
	assert.Equal(t, ErrTaskNotFound, err)

	// Test full queue
	task1 := &PipelineTask{ID: 1}
	task2 := &PipelineTask{ID: 2}
	task3 := &PipelineTask{ID: 3}

	q.AddToBack(task1)
	q.AddToBack(task2)
	err = q.AddToFront(task3)
	assert.Equal(t, ErrQueueFull, err)
}

func TestTaskLifecycle(t *testing.T) {
	q := NewTaskQueue(10)
	task := &PipelineTask{ID: 1}

	// Enqueue
	err := q.AddToBack(task)
	assert.NoError(t, err)
	assert.Equal(t, 1, q.PendingCount())

	// Dequeue to processing
	dequeued, err := q.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 0, q.PendingCount())
	assert.Equal(t, 1, q.ProcessingCount())

	// Complete
	q.CompleteTask(dequeued.ID)
	assert.Equal(t, 0, q.ProcessingCount())
}

func TestPendingMapConsistency(t *testing.T) {
	q := NewTaskQueue(10)
	task := &PipelineTask{ID: 1}

	q.AddToBack(task)
	assert.Equal(t, 1, q.pending.Len())
	assert.Equal(t, 1, len(q.pendingMap))

	// Remove and verify map is clean
	err := q.RemoveTaskByID(1)
	assert.NoError(t, err)
	assert.Equal(t, 0, q.pending.Len())
	assert.Equal(t, 0, len(q.pendingMap))

	// Test dequeue also cleans the map
	q.AddToBack(task)
	dequeued, err := q.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(q.pendingMap))

	q.CompleteTask(dequeued.ID)
}

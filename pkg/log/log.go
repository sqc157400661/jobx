package log

import (
	"bytes"
	"sync"
)

// LoggerSet manages a collection of log buffers
type LoggerSet struct {
	items map[int]*bytes.Buffer
	mu    sync.RWMutex
}

// NewLoggerSet creates a new LoggerSet instance
func NewLoggerSet() *LoggerSet {
	return &LoggerSet{
		items: make(map[int]*bytes.Buffer),
	}
}

// Add adds the item to the set.
func (s *LoggerSet) Add(item int) {
	buf := new(bytes.Buffer)
	s.add(item, buf)
}

// Add adds the item to the set.
func (s *LoggerSet) add(item int, buf *bytes.Buffer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[item] = buf
}

// AddOrGet retrieves or creates a buffer for index
func (s *LoggerSet) AddOrGet(item int) *bytes.Buffer {
	if s.Contains(item) {
		return s.Get(item)
	} else {
		buf := new(bytes.Buffer)
		s.add(item, buf)
		return buf
	}
}

func (s *LoggerSet) Get(item int) *bytes.Buffer {
	if s.Contains(item) {
		return s.items[item]
	}
	return nil
}

// Contains check if items (one or more) are present in the set.
func (s *LoggerSet) Contains(items ...int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, item := range items {
		if _, contains := s.items[item]; !contains {
			return false
		}
	}
	return true
}

// Empty returns true if set does not contain any elements.
func (s *LoggerSet) Empty() bool {
	return s.Size() == 0
}

// Size returns number of elements within the set.
func (s *LoggerSet) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

// ClearAll clears all values in the set.
func (s *LoggerSet) ClearAll() {
	s.items = make(map[int]*bytes.Buffer)
}

// Items all items in the set.
func (s *LoggerSet) Items() map[int]*bytes.Buffer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.items
}

func (s *LoggerSet) CleanupEmpty() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, buf := range s.items {
		if buf.Len() == 0 {
			delete(s.items, id)
		}
	}
}

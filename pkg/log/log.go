package log

import (
	"bytes"
	"sync"
)

func NewLoggerSet() *LoggerSet {
	return &LoggerSet{
		items: make(map[int]*bytes.Buffer),
	}
}

// Set holds elements in go's native map
type LoggerSet struct {
	items map[int]*bytes.Buffer
	lock  sync.RWMutex
	wg    sync.WaitGroup
}

// ReAdd adds the item to the set.
func (s *LoggerSet) ReAdd(item int, buf *bytes.Buffer) {
	s.add(item, buf)
}

// Add adds the item to the set.
func (s *LoggerSet) Add(item int) {
	buf := new(bytes.Buffer)
	s.add(item, buf)
}

// Add adds the item to the set.
func (s *LoggerSet) add(item int, buf *bytes.Buffer) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items[item] = buf
}

// Add adds the item to the set.
func (s *LoggerSet) AddOrGet(item int) *bytes.Buffer {
	if s.Contains(item) {
		return s.Get(item)
	} else {
		buf := new(bytes.Buffer)
		s.add(item, buf)
		return buf
	}
}

// Remove removes the items (one or more) from the set.
func (s *LoggerSet) Remove(items ...int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, item := range items {
		buf := s.Get(item)
		if buf != nil {
			buf.Reset()
			//buf = nil
		}
		delete(s.items, item)
	}
}

// Remove removes the items (one or more) from the set.
func (s *LoggerSet) Get(item int) *bytes.Buffer {
	if s.Contains(item) {
		return s.items[item]
	}
	return nil
}

// Contains check if items (one or more) are present in the set.
func (s *LoggerSet) Contains(items ...int) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
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
	return len(s.items)
}

// Clear clears all values in the set.
func (s *LoggerSet) Clear() {
	s.items = make(map[int]*bytes.Buffer)
}

// all items in the set.
func (s *LoggerSet) Items() map[int]*bytes.Buffer {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.items
}

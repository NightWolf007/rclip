// Package store represents a clipboard storage.
// The store implementation is threadsafe.
package store

import (
	"sync"
)

// Store represents store struct.
type Store struct {
	mu      sync.RWMutex
	list    *List
	maxSize uint
}

// New builds new store.
func New(maxSize uint) *Store {
	return &Store{
		list:    NewList(),
		maxSize: maxSize,
	}
}

// Get returns the latest element from store.
func (s *Store) Get() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.list.Last() == nil {
		return nil
	}

	return s.list.Last().Value
}

// GetAll returns a slice with all elements from store.
func (s *Store) GetAll() [][]byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data := make([][]byte, 0, s.list.Len())

	for e := s.list.Last(); e != nil; e = e.Prev() {
		data = append(data, e.Value)
	}

	return data
}

// Push pushes new data into the store.
func (s *Store) Push(data []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.list.Push(data)

	if s.list.Len() > s.maxSize {
		s.list.Pop()
	}
}

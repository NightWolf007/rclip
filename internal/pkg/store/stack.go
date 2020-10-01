package store

import "sync"

// Stack struct represents stack with fixed size.
type Stack struct {
	items [][]byte
	mutex sync.Mutex
}

// NewStack creates new stack with the given size.
func NewStack(size uint) *Stack {
	return &Stack{
		items: make([][]byte, size),
	}
}

// Last returns the latest value from the store.
func (s *Stack) Last() []byte {
	return s.Get(0)
}

// Get returns value from the stack by the given index.
func (s *Stack) Get(idx uint) []byte {
	if idx >= uint(len(s.items)) {
		return nil
	}

	return s.items[idx]
}

// GetAll returns all available values from the stack.
func (s *Stack) GetAll() [][]byte {
	vals := make([][]byte, 0, len(s.items))

	for i := 0; i < len(s.items) && s.items[i] != nil; i++ {
		vals = append(vals, s.items[i])
	}

	return vals
}

// Push inserts new element in the head of the stack.
func (s *Stack) Push(val []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i := len(s.items) - 1; i > 0; i-- {
		s.items[i] = s.items[i-1]
	}

	s.items[0] = val
}

// Pop deletes and returns the latest value from the stack.
func (s *Stack) Pop() []byte {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	val := s.items[0]

	for i := 0; i < len(s.items)-1; i++ {
		s.items[i] = s.items[i+1]
	}

	s.items[len(s.items)-1] = nil

	return val
}

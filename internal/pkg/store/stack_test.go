package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackLast(t *testing.T) {
	tests := []struct {
		name     string
		data     [][]byte
		expected []byte
	}{
		{
			name: "Simple",
			data: [][]byte{
				[]byte("1"),
				[]byte("2"),
				[]byte("3"),
			},
			expected: []byte("3"),
		},
		{
			name: "Overflowed",
			data: [][]byte{
				[]byte("1"),
				[]byte("2"),
				[]byte("3"),
				[]byte("4"),
				[]byte("5"),
				[]byte("6"),
				[]byte("7"),
			},
			expected: []byte("7"),
		},
		{
			name:     "EmptyList",
			data:     [][]byte{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := NewStack(5)

			for _, val := range tt.data {
				s.Push(val)
			}

			actual := s.Last()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestStackGet(t *testing.T) {
	tests := []struct {
		name     string
		data     [][]byte
		idx      uint
		expected []byte
	}{
		{
			name: "Simple",
			data: [][]byte{
				[]byte("1"),
				[]byte("2"),
				[]byte("3"),
			},
			idx:      1,
			expected: []byte("2"),
		},
		{
			name: "Overflowed",
			data: [][]byte{
				[]byte("1"),
				[]byte("2"),
				[]byte("3"),
				[]byte("4"),
				[]byte("5"),
				[]byte("6"),
				[]byte("7"),
			},
			idx:      1,
			expected: []byte("6"),
		},
		{
			name: "EmptyIndex",
			data: [][]byte{
				[]byte("1"),
				[]byte("2"),
			},
			idx:      2,
			expected: nil,
		},
		{
			name: "IndexOutOfRange",
			data: [][]byte{
				[]byte("1"),
				[]byte("2"),
				[]byte("3"),
			},
			idx:      10,
			expected: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := NewStack(5)

			for _, val := range tt.data {
				s.Push(val)
			}

			actual := s.Get(tt.idx)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestStackGetAll(t *testing.T) {
	tests := []struct {
		name     string
		data     [][]byte
		expected [][]byte
	}{
		{
			name: "Simple",
			data: [][]byte{
				[]byte("1"),
				[]byte("2"),
				[]byte("3"),
			},
			expected: [][]byte{
				[]byte("3"),
				[]byte("2"),
				[]byte("1"),
			},
		},
		{
			name: "Overflowed",
			data: [][]byte{
				[]byte("1"),
				[]byte("2"),
				[]byte("3"),
				[]byte("4"),
				[]byte("5"),
				[]byte("6"),
				[]byte("7"),
			},
			expected: [][]byte{
				[]byte("7"),
				[]byte("6"),
				[]byte("5"),
				[]byte("4"),
				[]byte("3"),
			},
		},
		{
			name:     "EmptyList",
			data:     [][]byte{},
			expected: [][]byte{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := NewStack(5)

			for _, val := range tt.data {
				s.Push(val)
			}

			actual := s.GetAll()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestStackPop(t *testing.T) {
	tests := []struct {
		name           string
		data           [][]byte
		expectedList   [][]byte
		expectedReturn []byte
	}{
		{
			name: "Simple",
			data: [][]byte{
				[]byte("1"),
				[]byte("2"),
				[]byte("3"),
			},
			expectedList: [][]byte{
				[]byte("2"),
				[]byte("1"),
			},
			expectedReturn: []byte("3"),
		},
		{
			name: "Overflowed",
			data: [][]byte{
				[]byte("1"),
				[]byte("2"),
				[]byte("3"),
				[]byte("4"),
				[]byte("5"),
				[]byte("6"),
				[]byte("7"),
			},
			expectedList: [][]byte{
				[]byte("6"),
				[]byte("5"),
				[]byte("4"),
				[]byte("3"),
			},
			expectedReturn: []byte("7"),
		},
		{
			name: "LastElement",
			data: [][]byte{
				[]byte("1"),
			},
			expectedList:   [][]byte{},
			expectedReturn: []byte("1"),
		},
		{
			name:           "EmptyList",
			data:           [][]byte{},
			expectedList:   [][]byte{},
			expectedReturn: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := NewStack(5)

			for _, val := range tt.data {
				s.Push(val)
			}

			actualReturn := s.Pop()
			assert.Equal(t, tt.expectedReturn, actualReturn)

			actualList := s.GetAll()
			assert.Equal(t, tt.expectedList, actualList)
		})
	}
}

package store_test

import (
	"testing"

	"github.com/NightWolf007/rclip/internal/pkg/store"
	"github.com/stretchr/testify/assert"
)

func TestStore_Get(t *testing.T) {
	tests := []struct {
		name     string
		data     [][]byte
		expected []byte
	}{
		{
			name:     "WhenStoreEmpty",
			data:     [][]byte{},
			expected: nil,
		},
		{
			name:     "WhenStoreNotEmpty",
			data:     [][]byte{{1}, {2}, {3}},
			expected: []byte{3},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := store.New(5)

			for _, val := range tt.data {
				s.Push(val)
			}

			actual := s.Get()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestStore_GetAll(t *testing.T) {
	tests := []struct {
		name     string
		data     [][]byte
		expected [][]byte
	}{
		{
			name:     "WhenStoreEmpty",
			data:     [][]byte{},
			expected: [][]byte{},
		},
		{
			name:     "WhenStoreNotEmpty",
			data:     [][]byte{{1}, {2}, {3}},
			expected: [][]byte{{3}, {2}, {1}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := store.New(5)

			for _, val := range tt.data {
				s.Push(val)
			}

			actual := s.GetAll()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestStore_GetPush(t *testing.T) {
	tests := []struct {
		name     string
		data     [][]byte
		value    []byte
		expected [][]byte
	}{
		{
			name:     "WhenStoreEmpty",
			data:     [][]byte{},
			value:    []byte{1},
			expected: [][]byte{{1}},
		},
		{
			name:     "WhenStoreNotEmpty",
			data:     [][]byte{{1}, {2}, {3}, {4}},
			value:    []byte{5},
			expected: [][]byte{{5}, {4}, {3}, {2}, {1}},
		},
		{
			name:     "WhenStoreOverflow",
			data:     [][]byte{{1}, {2}, {3}, {4}, {5}},
			value:    []byte{6},
			expected: [][]byte{{6}, {5}, {4}, {3}, {2}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := store.New(5)

			for _, val := range tt.data {
				s.Push(val)
			}

			s.Push(tt.value)

			actual := s.GetAll()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

package store_test

import (
	"testing"

	"github.com/NightWolf007/rclip/internal/pkg/store"
	"github.com/stretchr/testify/assert"
)

func TestList_FirstLast(t *testing.T) {
	tests := []struct {
		name          string
		data          [][]byte
		expectedFirst []byte
		expectedLast  []byte
	}{
		{
			name:          "WhenListEmpty",
			data:          [][]byte{},
			expectedFirst: nil,
			expectedLast:  nil,
		},
		{
			name:          "WhenOneElement",
			data:          [][]byte{{1}},
			expectedFirst: []byte{1},
			expectedLast:  []byte{1},
		},
		{
			name:          "WhenMultipleElements",
			data:          [][]byte{{1}, {2}, {3}},
			expectedFirst: []byte{1},
			expectedLast:  []byte{3},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			list := store.NewList()

			for _, val := range tt.data {
				list.Push(val)
			}

			actualFirst := list.First()
			actualLast := list.Last()

			if tt.expectedFirst == nil {
				assert.Nil(t, actualFirst)
			} else {
				assert.Equal(t, tt.expectedFirst, actualFirst.Value)
			}

			if tt.expectedLast == nil {
				assert.Nil(t, actualLast)
			} else {
				assert.Equal(t, tt.expectedLast, actualLast.Value)
			}
		})
	}
}

func TestList_Push(t *testing.T) {
	tests := []struct {
		name     string
		data     [][]byte
		value    []byte
		expected [][]byte
	}{
		{
			name:     "WhenListEmpty",
			data:     [][]byte{},
			value:    []byte{1},
			expected: [][]byte{{1}},
		},
		{
			name:  "WhenOneElement",
			data:  [][]byte{{1}},
			value: []byte{2},
			expected: [][]byte{
				{1},
				{2},
			},
		},
		{
			name:     "WhenMultipleElements",
			data:     [][]byte{{1}, {2}, {3}},
			value:    []byte{4},
			expected: [][]byte{{1}, {2}, {3}, {4}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			list := store.NewList()

			for _, val := range tt.data {
				list.Push(val)
			}

			list.Push(tt.value)

			assert.Equal(t, uint(len(tt.expected)), list.Len())

			i := 0
			for e := list.First(); e != nil; e = e.Next() {
				assert.Equal(t, tt.expected[i], e.Value)
				i++
			}
		})
	}
}

func TestList_Pop(t *testing.T) {
	tests := []struct {
		name           string
		data           [][]byte
		expectedReturn []byte
		expectedList   [][]byte
	}{
		{
			name:           "WhenListEmpty",
			data:           [][]byte{},
			expectedReturn: nil,
			expectedList:   [][]byte{},
		},
		{
			name:           "WhenOneElement",
			data:           [][]byte{{1}},
			expectedReturn: []byte{1},
			expectedList:   [][]byte{},
		},
		{
			name:           "WhenMultipleElements",
			data:           [][]byte{{1}, {2}, {3}},
			expectedReturn: []byte{1},
			expectedList:   [][]byte{{2}, {3}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			list := store.NewList()

			for _, val := range tt.data {
				list.Push(val)
			}

			actualReturn := list.Pop()
			if tt.expectedReturn == nil {
				assert.Nil(t, actualReturn)
			} else {
				assert.Equal(t, tt.expectedReturn, actualReturn.Value)
			}

			assert.Equal(t, uint(len(tt.expectedList)), list.Len())

			i := 0
			for e := list.First(); e != nil; e = e.Next() {
				assert.Equal(t, tt.expectedList[i], e.Value)
				i++
			}
		})
	}
}

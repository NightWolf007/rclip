package clipboard_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/NightWolf007/rclip/internal/pkg/clipboard"
)

type ClipboardMock struct {
	mock.Mock
}

func (m *ClipboardMock) Read() ([]byte, error) {
	args := m.Called()

	return args.Get(0).([]byte), args.Error(1)
}

func (m *ClipboardMock) Write(val []byte) error {
	args := m.Called(val)

	return args.Error(0)
}

func TestWatch(t *testing.T) {
	ctx, cancelFn := context.WithCancel(context.Background())
	cbMock := &ClipboardMock{}

	stream := clipboard.Watch(ctx)
	stream.Clipboard = cbMock

	timeAfterCounter := 0
	stream.TimeAfter = func(delay time.Duration) <-chan time.Time {
		timeAfterCounter++

		assert.Equal(t, time.Second, delay)

		ch := make(chan time.Time, 1)
		ch <- time.Now()

		return ch
	}

	cbMock.On("Read").Once().Return([]byte{1}, nil)
	cbMock.On("Read").Once().Return([]byte{1}, nil)
	cbMock.On("Read").Once().Return([]byte{1}, nil)
	cbMock.On("Read").Once().Return([]byte{1}, nil)
	cbMock.On("Read").Return([]byte{2}, nil)

	val, err := stream.Recv()
	assert.NoError(t, err)
	assert.Equal(t, []byte{1}, val)
	assert.Equal(t, 0, timeAfterCounter)

	timeAfterCounter = 0

	val, err = stream.Recv()
	assert.NoError(t, err)
	assert.Equal(t, []byte{2}, val)
	assert.Equal(t, 3, timeAfterCounter)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		_, err := stream.Recv()
		assert.Error(t, err)
		wg.Done()
	}()

	cancelFn()
	wg.Wait()

	cbMock.AssertExpectations(t)
}

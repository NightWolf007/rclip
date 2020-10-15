package grace_test

import (
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/NightWolf007/rclip/internal/app/grace"
	"github.com/stretchr/testify/assert"
)

// FnMock is mock helper for function calls.
type FnMock struct {
	mu    sync.Mutex
	calls []string
}

func NewFnMock() *FnMock {
	return &FnMock{
		calls: []string{},
	}
}

func (m *FnMock) Called(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.calls = append(m.calls, name)
}

func (m *FnMock) Calls() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.calls
}

func TestShutdown_Run(t *testing.T) {
	tests := []struct {
		name          string
		onWait        func(termCh chan<- os.Signal, timeCh chan<- time.Time)
		timeout       time.Duration
		expectedCalls []string
	}{
		{
			name:    "WhenShutdownIsDone",
			onWait:  func(termCh chan<- os.Signal, timeCh chan<- time.Time) {},
			timeout: 0,
			expectedCalls: []string{
				"SignalNotify",
				"OnShutdown",
				"Shutdown",
				"Wait",
				"OnDone",
			},
		},
		{
			name: "WhenSecondSIGTERMReceived",
			onWait: func(termCh chan<- os.Signal, timeCh chan<- time.Time) {
				termCh <- syscall.SIGTERM
			},
			timeout: 0,
			expectedCalls: []string{
				"SignalNotify",
				"OnShutdown",
				"Shutdown",
				"Wait",
				"OnForceQuit",
			},
		},
		{
			name: "WhenTimeout",
			onWait: func(termCh chan<- os.Signal, timeCh chan<- time.Time) {
				timeCh <- time.Now()
			},
			timeout: 1 * time.Nanosecond,
			expectedCalls: []string{
				"SignalNotify",
				"OnShutdown",
				"Shutdown",
				"Wait",
				"OnTimeout",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var termCh chan<- os.Signal
			timeCh := make(chan time.Time)

			fnMock := NewFnMock()

			g := &grace.Grace{
				Shutdown: func() {
					fnMock.Called("Shutdown")
				},
				Wait: func() {
					fnMock.Called("Wait")

					tt.onWait(termCh, timeCh)
				},
				Timeout: tt.timeout,

				OnShutdown: func() {
					fnMock.Called("OnShutdown")
				},
				OnDone: func() {
					fnMock.Called("OnDone")
				},
				OnForceQuit: func() {
					fnMock.Called("OnForceQuit")
				},
				OnTimeout: func() {
					fnMock.Called("OnTimeout")
				},

				SignalNotify: func(c chan<- os.Signal, sig ...os.Signal) {
					fnMock.Called("SignalNotify")

					assert.Equal(t, []os.Signal{syscall.SIGINT, syscall.SIGTERM}, sig)

					c <- syscall.SIGTERM
					termCh = c
				},

				TimeAfter: func(d time.Duration) <-chan time.Time {
					assert.Equal(t, tt.timeout, d)

					return timeCh
				},
			}

			g.Run()

			assert.Equal(t, tt.expectedCalls, fnMock.Calls())
		})
	}
}

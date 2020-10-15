package clipboard

import (
	"bytes"
	"context"
	"fmt"
	"time"
)

// Stream represents a clipboard changes stream.
type Stream struct {
	// RequestDelay is an interval between clipboard change checks.
	RequestDelay time.Duration
	// Base clipboard to use.
	// Overwrite it only for mocking purpose.
	Clipboard Clipboard

	TimeAfter func(delay time.Duration) <-chan time.Time

	ctx       context.Context
	prevValue []byte
}

// Watch creates a watch stream for clipboard changes.
func Watch(ctx context.Context) *Stream {
	return &Stream{
		RequestDelay: time.Second,
		Clipboard:    clipboardImpl{},
		TimeAfter:    time.After,
		ctx:          ctx,
	}
}

// Recv blocks thread until clipboard value change.
// Use context to cancel receiving.
func (s *Stream) Recv() ([]byte, error) {
	for {
		val, err := s.Clipboard.Read()
		if err != nil {
			return nil, fmt.Errorf("clipboard read: %w", err)
		}

		if !bytes.Equal(val, s.prevValue) {
			s.prevValue = val

			return val, nil
		}

		select {
		case <-s.TimeAfter(s.RequestDelay):
			continue
		case <-s.ctx.Done():
			return nil, s.ctx.Err()
		}
	}
}

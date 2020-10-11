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

	prevValue []byte
}

// Recv blocks thread until clipboard value change.
// Use context to cancel receiving.
func (s *Stream) Recv(ctx context.Context) ([]byte, error) {
	for {
		val, err := s.Clipboard.Read()
		if err != nil {
			return nil, fmt.Errorf("clipboard read: %w", err)
		}

		fmt.Printf("%v\n", val)

		if !bytes.Equal(val, s.prevValue) {
			s.prevValue = val

			return val, nil
		}

		select {
		case <-s.TimeAfter(s.RequestDelay):
			continue
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

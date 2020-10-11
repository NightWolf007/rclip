package pubsub_test

import (
	"sync"
	"testing"

	"github.com/NightWolf007/rclip/internal/pkg/pubsub"
	"github.com/stretchr/testify/assert"
)

func TestPubSub(t *testing.T) {
	value := []byte("hello")
	ps := pubsub.New()

	subWg := sync.WaitGroup{}
	wg := sync.WaitGroup{}

	for i := 0; i < 3; i++ {
		subWg.Add(1)
		wg.Add(1)

		go func() {
			defer wg.Done()

			sub := ps.Subscribe()

			subWg.Done()

			val := <-sub.C()
			assert.Equal(t, value, val)

			ps.Unsubscribe(sub)

			// Check that channel closed
			_, ok := <-sub.C()
			assert.False(t, ok)
		}()
	}

	subWg.Wait()

	ps.Publish(value)

	wg.Wait()
}

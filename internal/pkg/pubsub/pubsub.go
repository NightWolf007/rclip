// Package pubsub provides simple pub/sub mechanism.
// The pubsub implementation is threadsafe.
package pubsub

import (
	"sync"
)

// PubSub represents a pub/sub struct.
type PubSub struct {
	mu   sync.RWMutex
	subs map[*Subscription]struct{}
}

// New builds new PubSub struct.
func New() *PubSub {
	return &PubSub{
		subs: make(map[*Subscription]struct{}),
	}
}

// Subscribe creates new subscription in pub/sub.
func (ps *PubSub) Subscribe() *Subscription {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	sub := NewSubscription()
	ps.subs[sub] = struct{}{}

	return sub
}

// Unsubscribe unsubsccribes subscription from pub/sub.
func (ps *PubSub) Unsubscribe(sub *Subscription) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	delete(ps.subs, sub)
	sub.close()
}

// Publish sends value to all subscriptions.
func (ps *PubSub) Publish(value []byte) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for sub := range ps.subs {
		sub.publish(value)
	}
}

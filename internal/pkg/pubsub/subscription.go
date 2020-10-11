package pubsub

// Subscription represents a single pubsub subscription.
type Subscription struct {
	c chan []byte
}

// NewSubscription inits new Subscription struct.
func NewSubscription() *Subscription {
	return &Subscription{
		c: make(chan []byte),
	}
}

// C returns subscription channel.
func (s *Subscription) C() <-chan []byte {
	return s.c
}

// publish sends data to subscription.
func (s *Subscription) publish(value []byte) {
	s.c <- value
}

// close closes the subscription.
func (s *Subscription) close() {
	close(s.c)
}

package servers

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/NightWolf007/rclip/internal/pkg/api"
	"github.com/NightWolf007/rclip/internal/pkg/pubsub"
	"github.com/NightWolf007/rclip/internal/pkg/store"
)

// ClipboardServer represents the server for clipboard API.
type ClipboardServer struct {
	mu     sync.Mutex
	store  *store.Store
	pubsub *pubsub.PubSub
}

// NewClipboardServer builds new ClipboardServer.
func NewClipboardServer(storeSize uint) *ClipboardServer {
	return &ClipboardServer{
		store:  store.New(storeSize),
		pubsub: pubsub.New(),
	}
}

// Get returns the latest value from clipboard.
func (s *ClipboardServer) Get(ctx context.Context, in *api.GetRequest) (*api.GetResponse, error) {
	val := s.store.Get()

	return &api.GetResponse{Value: val}, nil
}

// Hist returns all clipboard history.
func (s *ClipboardServer) Hist(ctx context.Context, in *api.HistRequest) (*api.HistResponse, error) {
	vals := s.store.GetAll()

	return &api.HistResponse{Values: vals}, nil
}

// Push writes new clipboard value.
func (s *ClipboardServer) Push(ctx context.Context, in *api.PushRequest) (*api.PushResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if in.Value != nil && !bytes.Equal(in.Value, s.store.Get()) {
		s.store.Push(in.Value)
		s.pubsub.Publish(in.Value)
	}

	return &api.PushResponse{}, nil
}

// Subscribe allows to subscribe on server clipboard changes.
func (s *ClipboardServer) Subscribe(in *api.SubscribeRequest, stream api.ClipboardAPI_SubscribeServer) error {
	sub := s.pubsub.Subscribe()
	defer s.pubsub.Unsubscribe(sub)

	err := stream.Send(&api.SubscribeResponse{
		Value: s.store.Get(),
	})
	if err != nil {
		return fmt.Errorf("stream send: %w", err)
	}

	for {
		select {
		case val := <-sub.C():
			err := stream.Send(&api.SubscribeResponse{
				Value: val,
			})
			if err != nil {
				return fmt.Errorf("stream send: %w", err)
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}

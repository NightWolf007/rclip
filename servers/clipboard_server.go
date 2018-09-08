package servers

import (
	"context"
	"sync"

	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/ksuid"

	"github.com/NightWolf007/rclip/pb"
)

type Subscriber struct {
	uid string
	C   chan []byte
}

func NewSubscriber() Subscriber {
	return Subscriber{
		uid: ksuid.New().String(),
		C:   make(chan []byte),
	}
}

type ClipboardServer struct {
	buffer      []byte
	subscribers map[string]Subscriber
	lock        *sync.Mutex
}

func NewClipboardServer() *ClipboardServer {
	return &ClipboardServer{
		buffer:      []byte{},
		subscribers: map[string]Subscriber{},
		lock:        &sync.Mutex{},
	}
}

func (s *ClipboardServer) Push(ctx context.Context, params *pb.PushRequest) (*pb.Clip, error) {
	log.Debug().Bytes("data", params.Data).Msg("Received")

	s.buffer = params.Data
	s.broadcast(s.buffer)

	return &pb.Clip{
		Data: s.buffer,
	}, nil
}

func (s *ClipboardServer) Get(ctx context.Context, params *pb.GetRequest) (*pb.Clip, error) {
	log.Debug().Bytes("data", s.buffer).Msg("Sending")

	return &pb.Clip{
		Data: s.buffer,
	}, nil
}

func (s *ClipboardServer) Subscribe(params *pb.SubscribeRequest, stream pb.Clipboard_SubscribeServer) error {
	sub := s.subscribe()
	defer s.unsubscribe(sub)

	stream.Send(&pb.Clip{
		Data: s.buffer,
	})

	for {
		select {
		case data := <-sub.C:
			stream.Send(&pb.Clip{
				Data: data,
			})
		case <-stream.Context().Done():
			return nil
		}
	}

	return nil
}

func (s *ClipboardServer) subscribe() Subscriber {
	sub := NewSubscriber()

	s.lock.Lock()
	defer s.lock.Unlock()

	s.subscribers[sub.uid] = sub

	log.Debug().Str("SubID", sub.uid).Msg("Subsribe")

	return sub
}

func (s *ClipboardServer) unsubscribe(sub Subscriber) {
	log.Debug().Str("SubID", sub.uid).Msg("Unsubsribe")

	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.subscribers, sub.uid)
}

func (s *ClipboardServer) broadcast(data []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, sub := range s.subscribers {
		sub.C <- data
	}
}

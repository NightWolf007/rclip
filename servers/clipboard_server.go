package servers

import (
	"context"

	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
	// "github.com/rs/zerolog/log"

	"github.com/NightWolf007/rclipd/pb"
)

type ClipboardServer struct {
	buffer []string
}

func NewClipboardServer() *ClipboardServer {
	return &ClipboardServer{
		buffer: []string{},
	}
}

func (s *ClipboardServer) Push(
	ctx context.Context, params *pb.PushRequest,
) (*pb.PushResponse, error) {
	s.buffer = append(s.buffer, params.Data)
	return &pb.PushResponse{}, nil
}

func (s *ClipboardServer) Get(
	ctx context.Context, params *pb.GetRequest,
) (*pb.GetResponse, error) {
	var data string
	if len(s.buffer) > 0 {
		data = s.buffer[len(s.buffer)-1]
	}
	return &pb.GetResponse{
		Data: data,
	}, nil
}

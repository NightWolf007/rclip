package servers

import (
	"context"

	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
	"github.com/rs/zerolog/log"

	"github.com/NightWolf007/rclip/pb"
)

type ClipboardServer struct {
	buffer string
}

func NewClipboardServer() *ClipboardServer {
	return &ClipboardServer{
		buffer: "",
	}
}

func (s *ClipboardServer) Push(
	ctx context.Context, params *pb.PushRequest,
) (*pb.PushResponse, error) {
	log.Debug().Msgf("Received data: %s", params.Data)
	s.buffer = params.Data
	return &pb.PushResponse{}, nil
}

func (s *ClipboardServer) Get(
	ctx context.Context, params *pb.GetRequest,
) (*pb.GetResponse, error) {
	log.Debug().Msgf("Sending data: %s", s.buffer)
	return &pb.GetResponse{
		Data: s.buffer,
	}, nil
}

// Package syncer provides methods to sync system clipboard with server clipboard.
package syncer

import (
	"context"
	"fmt"

	"github.com/NightWolf007/rclip/internal/pkg/api"
	"github.com/NightWolf007/rclip/internal/pkg/clipboard"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// Syncer struct provides clipboard sync methods.
type Syncer struct {
	addr   string
	logger zerolog.Logger
}

// New builds new syncer.
func New(addr string, logger zerolog.Logger) *Syncer {
	return &Syncer{
		addr:   addr,
		logger: logger,
	}
}

// RemoteToLocal starts sync loop from remote clipboard to local clipboard.
func (s *Syncer) RemoteToLocal(ctx context.Context) error {
	conn, err := grpc.Dial(s.addr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("grpc dial: %w", err)
	}

	defer conn.Close()

	client := api.NewClipboardAPIClient(conn)

	stream, err := client.Subscribe(ctx, &api.SubscribeRequest{})
	if err != nil {
		return fmt.Errorf("client subscribe: %w", err)
	}

	return s.consumeLoop(
		ctx,
		NewRemoteInputStream(stream),
		NewLocalOutputStream(),
	)
}

// LocalToRemote starts sync loop from local clipboard to remote clipboard.
func (s *Syncer) LocalToRemote(ctx context.Context) error {
	conn, err := grpc.Dial(s.addr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("grpc dial: %w", err)
	}

	defer conn.Close()

	return s.consumeLoop(
		ctx,
		NewLocalInputStream(clipboard.Watch(ctx)),
		NewRemoteOutputStream(ctx, api.NewClipboardAPIClient(conn)),
	)
}

func (s *Syncer) consumeLoop(ctx context.Context, in InputStream, out OutputStream) error {
	s.logger.Debug().
		Msg("Starting consuming loop")

	for {
		err := s.consume(ctx, in, out)
		if err != nil {
			return fmt.Errorf("consume: %w", err)
		}
	}
}

func (s *Syncer) consume(ctx context.Context, in InputStream, out OutputStream) error {
	val, err := in.Recv()
	if err != nil {
		if ctx.Err() != nil {
			return nil
		}

		return fmt.Errorf("stream recv: %w", err)
	}

	s.logger.Debug().
		Bytes("value", val).
		Msg("Received new value from input stream")

	if val != nil {
		err := out.Send(val)
		if err != nil {
			return fmt.Errorf("stream send: %w", err)
		}

		s.logger.Debug().
			Bytes("value", val).
			Msg("New value has been sent to the output stream")
	}

	return nil
}

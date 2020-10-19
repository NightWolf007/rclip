// Package syncer provides methods to sync system clipboard with server clipboard.
package syncer

import (
	"context"
	"fmt"

	"github.com/NightWolf007/rclip/internal/app/client"
	"github.com/NightWolf007/rclip/internal/pkg/api"
	"github.com/NightWolf007/rclip/internal/pkg/clipboard"
	"github.com/rs/zerolog"
)

// Syncer struct provides clipboard sync methods.
type Syncer struct {
	logger zerolog.Logger
}

// New builds new syncer.
func New(logger zerolog.Logger) *Syncer {
	return &Syncer{
		logger: logger,
	}
}

// RemoteToLocal starts sync loop from remote clipboard to local clipboard.
func (s *Syncer) RemoteToLocal(ctx context.Context) error {
	cli, err := client.DialWithLogger(s.logger)
	if err != nil {
		return fmt.Errorf("client dial: %w", err)
	}

	defer cli.Close()

	stream, err := cli.Subscribe(ctx, &api.SubscribeRequest{})
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
	cli, err := client.DialWithLogger(s.logger)
	if err != nil {
		return fmt.Errorf("client dial: %w", err)
	}

	defer cli.Close()

	return s.consumeLoop(
		ctx,
		NewLocalInputStream(clipboard.Watch(ctx)),
		NewRemoteOutputStream(ctx, cli),
	)
}

func (s *Syncer) consumeLoop(ctx context.Context, in InputStream, out OutputStream) (err error) {
	log := s.logger

	log.Debug().
		Msg("Starting consuming loop")

	for {
		err = s.consume(in, out)
		if err != nil {
			if ctx.Err() != nil {
				log.Debug().Err(err).Msg("Context is done")
				err = nil
			}

			break
		}
	}

	log.Debug().
		Err(err).
		Msg("Exiting consuming loop")

	return
}

func (s *Syncer) consume(in InputStream, out OutputStream) error {
	log := s.logger

	log.Debug().
		Msg("Receiving value from the input stream")

	val, err := in.Recv()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Receive on input stream failed")

		return fmt.Errorf("stream recv: %w", err)
	}

	log = log.With().Bytes("value", val).Logger()

	log.Debug().
		Bytes("value", val).
		Msg("Received new value from input stream")

	if val != nil {
		log.Debug().
			Msg("Sending value to the output stream")

		err := out.Send(val)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to write value to the output stream")

			return fmt.Errorf("stream send: %w", err)
		}

		log.Debug().
			Msg("Value sent to the output stream")
	}

	return nil
}

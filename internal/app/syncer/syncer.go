// Package syncer provides methods to sync system clipboard with server clipboard.
package syncer

import (
	"bytes"
	"context"
	"fmt"

	"github.com/NightWolf007/rclip/internal/pkg/api"
	"github.com/NightWolf007/rclip/internal/pkg/clipboard"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Syncer struct {
	serverAddr string
}

func RemoteToLocal(ctx context.Context) error {
	serverAddr := "localhost:9889"

	for {

		err = syncRemoteToLocal(ctx)
	}
}

func syncRemoteToLocal(ctx context.Context) error {
	serverAddr := "localhost:9889"

	log := log.With().Str("addr", serverAddr).Logger()

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to connect to the server")

		return fmt.Errorf("grpc dial: %w", err)
	}

	defer conn.Close()

	client := api.NewClipboardAPIClient(conn)

	stream, err := client.Subscribe(ctx, &api.SubscribeRequest{})
	if err != nil {
		log.Error().
			Err(err).
			Str("method", "Subscribe").
			Msg("Failed to execute RPC method")

		return fmt.Errorf("clipboardapi client subscribe: %w")
	}

	for {
		select {}
		resp, err := stream.Recv()
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to receive value from stream")

			return fmt.Errorf("stream recv: %w", err)
		}

		log = log.With().Bytes("value", resp.Value).Logger()

		log.Debug().
			Msg("Received new value from stream")

		updateClipboardValue(log, resp.Value)
	}
}

func updateClipboardValue(log zerolog.Logger, val []byte) {
	currentVal, err := clipboard.Read()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to read clipboard current value")
	}

	log = log.With().Bytes("currentValue", val).Logger()

	if val != nil && !bytes.Equal(val, currentVal) {
		err := clipboard.Write(val)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to write value to clipboard")

			return
		}

		log.Debug().
			Msg("The value written to the clipboard")
	}
}

func daemonListenClipboard(ctx context.Context) error {
	stream := clipboard.Watch()

	conn, err := grpc.Dial(pasteListenAddr, grpc.WithInsecure())
	if err != nil {
		log.Error().
			Err(err).
			Str("addr", copyListenAddr).
			Msg("Failed to connect to the server")

		return fmt.Errorf("grpc dial: %w", err)
	}

	defer conn.Close()

	client := api.NewClipboardAPIClient(conn)

	for {
		val, err := stream.Recv(ctx)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to receive value from clipboard stream")

			continue
		}

		_, err = client.Push(ctx, &api.PushRequest{Value: val})
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to push clipboard value to the server")

			return fmt.Errorf("client push: %w", err)
		}
	}
}

// Package client provides RClip client helpers.
package client

import (
	"fmt"

	"github.com/NightWolf007/rclip/internal/app/client/interceptors"
	"github.com/NightWolf007/rclip/internal/pkg/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// Client represents commpon ClipboardAPIClient wrapper.
type Client struct {
	api.ClipboardAPIClient
	logger zerolog.Logger
	conn   *grpc.ClientConn
}

// Dial creates new client with established connection to the server.
func Dial() (*Client, error) {
	return DialWithLogger(log.Logger)
}

// DialWithLogger creates new client with established connection to the server.
// It uses the given logger.
func DialWithLogger(logger zerolog.Logger) (*Client, error) {
	target := viper.GetString("client.target")
	log := logger.With().Str("target", target).Logger()

	log.Debug().Msg("Dialing server")

	conn, err := grpc.Dial(
		target,
		grpc.WithInsecure(),
		grpc.WithChainStreamInterceptor(
			interceptors.StreamLogInterceptor(log),
		),
		grpc.WithChainUnaryInterceptor(
			interceptors.UnaryLogInterceptor(log),
		),
	)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to establish connection")

		return nil, fmt.Errorf("grpc dial: %w", err)
	}

	log.Debug().Msg("Connection established successfully")

	return &Client{
		ClipboardAPIClient: api.NewClipboardAPIClient(conn),
		conn:               conn,
		logger:             log,
	}, nil
}

// Close closes client connection.
func (c *Client) Close() error {
	log := c.logger

	log.Debug().Msg("Closing connection")

	err := c.conn.Close()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to close connection")

		return fmt.Errorf("grpc connection close: %w", err)
	}

	log.Debug().Msg("Connection successfully closed")

	return nil
}

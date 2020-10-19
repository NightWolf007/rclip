package interceptors

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// UnaryLogInterceptor returns a new unary server interceptor that logs incomming gRPC calls.
func UnaryLogInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log := logger.With().
			Str("method", info.FullMethod).
			Interface("request", req).
			Logger()

		log.Debug().Msg("Unary RPC received")

		start := time.Now()
		reply, err := handler(ctx, req)
		dur := time.Since(start)

		log = log.With().
			Dur("duration", dur).
			Interface("reply", reply).
			Logger()

		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to process unary RPC")

			return reply, err
		}

		log.Debug().
			Err(err).
			Msg("Unary RPC call successfully processed")

		return reply, nil
	}
}

// StreamLogInterceptor returns a new streaming client interceptor that logs the execution of gRPC calls.
func StreamLogInterceptor(logger zerolog.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		var streamType string

		switch {
		case info.IsClientStream && info.IsServerStream:
			streamType = "bidirectional"
		case info.IsClientStream:
			streamType = "client"
		case info.IsServerStream:
			streamType = "server"
		}

		log := logger.With().
			Str("method", info.FullMethod).
			Str("type", streamType).
			Logger()

		log.Debug().
			Msg("Streaming RPC received")

		start := time.Now()
		err := handler(srv, newLogInterceptorStreamWrapper(ss, log))
		dur := time.Since(start)

		log = log.With().
			Dur("duration", dur).
			Logger()

		if err != nil {
			log.Error().
				Err(err).
				Msg("RPC stream failed")

			return err
		}

		log.Debug().
			Msg("RPC stream finished successfully")

		return nil
	}
}

type logInterceptorStreamWrapper struct {
	grpc.ServerStream

	logger zerolog.Logger
}

func newLogInterceptorStreamWrapper(stream grpc.ServerStream, logger zerolog.Logger) *logInterceptorStreamWrapper {
	return &logInterceptorStreamWrapper{
		ServerStream: stream,
		logger:       logger,
	}
}

func (s *logInterceptorStreamWrapper) SendMsg(m interface{}) error {
	log := s.logger.With().Interface("message", m).Logger()

	log.Debug().
		Msg("Sending message to stream")

	err := s.ServerStream.SendMsg(m)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to send message to stream")

		return err
	}

	log.Debug().
		Msg("Message successfully sent to stream")

	return nil
}

func (s *logInterceptorStreamWrapper) RecvMsg(m interface{}) error {
	log := s.logger

	log.Debug().
		Msg("Receiving message from stream")

	err := s.ServerStream.RecvMsg(m)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to receive message from stream")

		return err
	}

	log.Debug().
		Interface("message", m).
		Msg("Received new message from stream")

	return nil
}

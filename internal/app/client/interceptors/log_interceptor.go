package interceptors

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// UnaryLogInterceptor returns a new unary client interceptor that logs the execution of gRPC calls.
func UnaryLogInterceptor(logger zerolog.Logger) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log := logger.With().
			Str("method", method).
			Interface("request", req).
			Interface("options", opts).
			Logger()

		log.Debug().Msg("Invoking RPC method")

		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		dur := time.Since(start)

		log = log.With().
			Dur("duration", dur).
			Interface("reply", reply).
			Logger()

		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to invoke RPC method")

			return err
		}

		log.Debug().
			Msg("RPC method invoked successfully")

		return err
	}
}

// StreamLogInterceptor returns a new streaming client interceptor that logs the execution of gRPC calls.
func StreamLogInterceptor(logger zerolog.Logger) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		log := logger.With().
			Str("method", method).
			Logger()

		log.Debug().
			Msg("Openning RPC stream")

		start := time.Now()
		stream, err := streamer(ctx, desc, cc, method, opts...)
		dur := time.Since(start)

		log = log.With().
			Dur("duration", dur).
			Logger()

		if err != nil {
			log.Error().
				Err(err).
				Msg("RPC stream failed")

			return stream, err
		}

		log.Debug().
			Msg("RPC stream finished successfully")

		return newLogInterceptorStreamWrapper(stream, log), nil
	}
}

type logInterceptorStreamWrapper struct {
	grpc.ClientStream

	logger zerolog.Logger
}

func newLogInterceptorStreamWrapper(stream grpc.ClientStream, logger zerolog.Logger) *logInterceptorStreamWrapper {
	return &logInterceptorStreamWrapper{
		ClientStream: stream,
		logger:       logger,
	}
}

func (s *logInterceptorStreamWrapper) SendMsg(m interface{}) error {
	log := s.logger.With().Interface("message", m).Logger()

	log.Debug().
		Msg("Sending message to stream")

	err := s.ClientStream.SendMsg(m)
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

	err := s.ClientStream.RecvMsg(m)
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

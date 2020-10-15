package syncer

import (
	"context"

	"github.com/NightWolf007/rclip/internal/pkg/api"
	"github.com/NightWolf007/rclip/internal/pkg/clipboard"
)

// InputStream is an abstract input stream.
type InputStream interface {
	Recv() ([]byte, error)
}

// OutputStream is an abstract output stream.
type OutputStream interface {
	Send([]byte) error
}

// RemoteInputStream represents an input stream from remote clipboard.
type RemoteInputStream struct {
	stream api.ClipboardAPI_SubscribeClient
}

// NewRemoteInputStream creates new RemoteInputStream.
func NewRemoteInputStream(stream api.ClipboardAPI_SubscribeClient) RemoteInputStream {
	return RemoteInputStream{stream: stream}
}

// Recv waits until new value from remote clipboard is received.
func (s RemoteInputStream) Recv() ([]byte, error) {
	resp, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}

// RemoteOutputStream represents an output stream to remote clipboard.
type RemoteOutputStream struct {
	ctx    context.Context
	client api.ClipboardAPIClient
}

// NewRemoteOutputStream create new RemoteOutputStream.
func NewRemoteOutputStream(ctx context.Context, client api.ClipboardAPIClient) RemoteOutputStream {
	return RemoteOutputStream{ctx: ctx, client: client}
}

// Send sends value to remote clipboard.
func (s RemoteOutputStream) Send(val []byte) error {
	_, err := s.client.Push(s.ctx, &api.PushRequest{Value: val})

	return err
}

// LocalInputStream represents an input stream from system clipboard.
type LocalInputStream struct {
	stream *clipboard.Stream
}

// NewLocalInputStream creates new LocalInputStream.
func NewLocalInputStream(stream *clipboard.Stream) LocalInputStream {
	return LocalInputStream{stream: stream}
}

// Recv waits until new value appears in the system clipboard.
func (s LocalInputStream) Recv() ([]byte, error) {
	return s.stream.Recv()
}

// LocalOutputStream represents an output stream to system clipboard.
type LocalOutputStream struct{}

// NewLocalOutputStream creates new LocalOutputStream.
func NewLocalOutputStream() LocalOutputStream {
	return LocalOutputStream{}
}

// Send updates value of the system clipboard.
func (s LocalOutputStream) Send(val []byte) error {
	return clipboard.Write(val)
}

package servers_test

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/NightWolf007/rclip/internal/app/servers"
	"github.com/NightWolf007/rclip/internal/pkg/api"
)

func TestClipboardServer_Get(t *testing.T) {
	ctx := context.Background()

	srv, lis := startClipboardServer(t)
	defer srv.Stop()

	conn := dialClipboardServer(t, lis)
	defer conn.Close()

	client := api.NewClipboardAPIClient(conn)

	t.Run("WhenStoreEmpty", func(t *testing.T) {
		resp, err := client.Get(ctx, &api.GetRequest{})
		assert.NoError(t, err)
		assert.Nil(t, resp.Value)
	})

	t.Run("WhenStoreHasOneElement", func(t *testing.T) {
		_, err := client.Push(ctx, &api.PushRequest{Value: []byte{1}})
		require.NoError(t, err)

		resp, err := client.Get(ctx, &api.GetRequest{})
		assert.NoError(t, err)
		assert.Equal(t, []byte{1}, resp.Value)
	})

	t.Run("WhenStoreHasMultipleElements", func(t *testing.T) {
		_, err := client.Push(ctx, &api.PushRequest{Value: []byte{2}})
		require.NoError(t, err)

		resp, err := client.Get(ctx, &api.GetRequest{})
		assert.NoError(t, err)
		assert.Equal(t, []byte{2}, resp.Value)
	})
}

func TestClipboardServer_Hist(t *testing.T) {
	ctx := context.Background()

	srv, lis := startClipboardServer(t)
	defer srv.Stop()

	conn := dialClipboardServer(t, lis)
	defer conn.Close()

	client := api.NewClipboardAPIClient(conn)

	t.Run("WhenStoreEmpty", func(t *testing.T) {
		resp, err := client.Hist(ctx, &api.HistRequest{})
		assert.NoError(t, err)
		assert.Empty(t, resp.Values)
	})

	t.Run("WhenStoreHasOneElement", func(t *testing.T) {
		_, err := client.Push(ctx, &api.PushRequest{Value: []byte{1}})
		require.NoError(t, err)

		resp, err := client.Hist(ctx, &api.HistRequest{})
		assert.NoError(t, err)
		assert.Equal(t, [][]byte{{1}}, resp.Values)
	})

	t.Run("WhenStoreHasMultipleElements", func(t *testing.T) {
		_, err := client.Push(ctx, &api.PushRequest{Value: []byte{2}})
		require.NoError(t, err)

		resp, err := client.Hist(ctx, &api.HistRequest{})
		assert.NoError(t, err)
		assert.Equal(t, [][]byte{{2}, {1}}, resp.Values)
	})
}

func TestClipboardServer_Push(t *testing.T) {
	ctx := context.Background()

	srv, lis := startClipboardServer(t)
	defer srv.Stop()

	conn := dialClipboardServer(t, lis)
	defer conn.Close()

	client := api.NewClipboardAPIClient(conn)

	t.Run("WhenStoreEmpty", func(t *testing.T) {
		_, err := client.Push(ctx, &api.PushRequest{Value: []byte{1}})
		assert.NoError(t, err)

		resp, err := client.Hist(ctx, &api.HistRequest{})
		require.NoError(t, err)
		assert.Equal(t, [][]byte{{1}}, resp.Values)
	})

	t.Run("WhenStoreHasElements", func(t *testing.T) {
		_, err := client.Push(ctx, &api.PushRequest{Value: []byte{2}})
		assert.NoError(t, err)

		resp, err := client.Hist(ctx, &api.HistRequest{})
		require.NoError(t, err)
		assert.Equal(t, [][]byte{{2}, {1}}, resp.Values)
	})

	t.Run("WhenDuplicateValue", func(t *testing.T) {
		_, err := client.Push(ctx, &api.PushRequest{Value: []byte{2}})
		assert.NoError(t, err)

		resp, err := client.Hist(ctx, &api.HistRequest{})
		require.NoError(t, err)
		assert.Equal(t, [][]byte{{2}, {1}}, resp.Values)
	})

	t.Run("WhenStoreOverflow", func(t *testing.T) {
		for i := byte(3); i < 7; i++ {
			_, err := client.Push(ctx, &api.PushRequest{Value: []byte{i}})
			require.NoError(t, err)
		}

		resp, err := client.Hist(ctx, &api.HistRequest{})
		assert.NoError(t, err)
		assert.Equal(t, [][]byte{{6}, {5}, {4}, {3}, {2}}, resp.Values)
	})
}

func TestClipboardServer_Subscribe(t *testing.T) {
	ctx := context.Background()

	srv, lis := startClipboardServer(t)
	defer srv.Stop()

	conn := dialClipboardServer(t, lis)
	defer conn.Close()

	client := api.NewClipboardAPIClient(conn)

	var (
		err     error
		stream1 api.ClipboardAPI_SubscribeClient
		stream2 api.ClipboardAPI_SubscribeClient
	)

	t.Run("WhenStoreEmpty", func(t *testing.T) {
		stream1, err = client.Subscribe(ctx, &api.SubscribeRequest{})
		assert.NoError(t, err)

		resp1, err := stream1.Recv()
		assert.NoError(t, err)
		assert.Nil(t, err, resp1.Value)

		_, err = client.Push(ctx, &api.PushRequest{Value: []byte{1}})
		require.NoError(t, err)

		resp1, err = stream1.Recv()
		assert.NoError(t, err)
		assert.Equal(t, []byte{1}, resp1.Value)
	})

	t.Run("WhenStoreHasElements", func(t *testing.T) {
		stream2, err = client.Subscribe(ctx, &api.SubscribeRequest{})
		assert.NoError(t, err)

		resp2, err := stream2.Recv()
		assert.NoError(t, err)
		assert.Nil(t, err, resp2.Value)

		_, err = client.Push(ctx, &api.PushRequest{Value: []byte{2}})
		require.NoError(t, err)

		resp1, err := stream1.Recv()
		assert.NoError(t, err)
		assert.Equal(t, []byte{2}, resp1.Value)

		resp2, err = stream2.Recv()
		assert.NoError(t, err)
		assert.Equal(t, []byte{2}, resp2.Value)
	})

	assert.NoError(t, stream1.CloseSend())
	assert.NoError(t, stream2.CloseSend())
}

func dialClipboardServer(t *testing.T, lis *bufconn.Listener) *grpc.ClientConn {
	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithInsecure(),
	)
	require.NoError(t, err, "failed to dial to GRPC server")

	return conn
}

func startClipboardServer(t *testing.T) (*grpc.Server, *bufconn.Listener) {
	bufferSize := 1024 * 1024
	lis := bufconn.Listen(bufferSize)
	srv := grpc.NewServer()

	clipboardServer := servers.NewClipboardServer(5)
	api.RegisterClipboardAPIServer(srv, clipboardServer)

	go func() {
		err := srv.Serve(lis)
		require.NoError(t, err, "failed to start GRPC server")
	}()

	return srv, lis
}

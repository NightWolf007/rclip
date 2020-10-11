package adapter

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisAdapter_Push(t *testing.T) {
	ctx := context.Background()
	data := "hello world"

	adapter := NewRedisAdapter(&redis.Options{
		Addr: "redis:6379",
		DB:   1,
	})

	require.NoError(t, adapter.Open(ctx))

	defer func() {
		require.NoError(t, adapter.Close(ctx))
	}()

	err := adapter.rdb.FlushDB(ctx).Err()
	require.NoError(t, err)

	err = adapter.Push(ctx, []byte(data))
	assert.NoError(t, err)

	actualData, err := adapter.rdb.LRange(ctx, RedisHistoryList, 0, 99).Result()
	require.NoError(t, err)
	assert.Equal(t, []string{data}, actualData)
}

// func TestRedisAdapter_Get(t *testing.T) {
// 	ctx := context.Background()
// 	data := "hello world"
//
// 	adapter := NewRedisAdapter(&redis.Options{
// 		Addr: "redis:6379",
// 		DB:   1,
// 	})
//
// 	require.NoError(t, adapter.Open(ctx))
//
// 	defer func() {
// 		require.NoError(t, adapter.Close(ctx))
// 	}()
//
// 	err := adapter.rdb.FlushDB(ctx).Err()
// 	require.NoError(t, err)
//
// 	err = adapter.Push(ctx, []byte(data))
// 	assert.NoError(t, err)
//
// 	actualData, err := adapter.rdb.LRange(ctx, RedisHistoryList, 0, 99).Result()
// 	require.NoError(t, err)
// 	assert.Equal(t, []string{data}, actualData)
// }

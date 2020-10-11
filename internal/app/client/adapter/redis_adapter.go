package adapter

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

const (
	// RedisHistoryList represents Redis key name of the history list.
	RedisHistoryList = "rclip_history"
	// RedisPubSubQueue represents Redis key name of the pub/sub queue.
	RedisPubSubQueue = "rclip_queue"
)

// RedisAdapter using Redis as server for RClip.
type RedisAdapter struct {
	opts *redis.Options
	rdb  *redis.Client
}

// NewRedisAdapter builds new RedisAdapter.
func NewRedisAdapter(opts *redis.Options) *RedisAdapter {
	return &RedisAdapter{opts: opts}
}

// Close closes adapter.
// It closes connection to Redis database.
func (a RedisAdapter) Close(ctx context.Context) error {
	err := a.rdb.Close()
	if err != nil {
		return fmt.Errorf("redis close: %w", err)
	}

	return nil
}

// Get returns the latest clipboard data.
func (a RedisAdapter) Get(ctx context.Context) ([]byte, error) {
	data, err := a.rdb.LIndex(ctx, RedisHistoryList, 0).Bytes()
	if err != nil {
		return nil, fmt.Errorf("redis lindex: %w", err)
	}

	return data, nil
}

// GetAll returns all clipboard history available.
func (a RedisAdapter) GetAll(ctx context.Context) ([][]byte, error) {
	dataStr, err := a.rdb.LRange(ctx, RedisHistoryList, 0, 99).Result()
	if err != nil {
		return nil, fmt.Errorf("redis lrange: %w", err)
	}

	data := make([][]byte, len(dataStr))
	for i := range dataStr {
		data[i] = []byte(dataStr[i])
	}

	return data, nil
}

// Open starts adapter.
// It opens connection to Redis database.
func (a *RedisAdapter) Open(ctx context.Context) error {
	a.rdb = redis.NewClient(a.opts)

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}

	return nil
}

// Push pushes data to the Redis list.
func (a RedisAdapter) Push(ctx context.Context, data []byte) error {
	pipe := a.rdb.TxPipeline()
	pipe.LPush(ctx, RedisHistoryList, data)
	pipe.LTrim(ctx, RedisHistoryList, 0, 99)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("redis pipeline exec: %w", err)
	}

	return nil
}

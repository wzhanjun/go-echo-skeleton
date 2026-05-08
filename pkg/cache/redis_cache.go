package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() *RedisCache {
	return &RedisCache{client: GetRedisConn()}
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	return c.client.Get(ctx, key).Bytes()
}

func (c *RedisCache) Set(ctx context.Context, key string, data []byte, ttl int) error {
	return c.client.Set(ctx, key, data, time.Duration(ttl)*time.Second).Err()
}

func (c *RedisCache) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

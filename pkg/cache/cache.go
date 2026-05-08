package cache

import "context"

// Cache 是通用缓存后端接口。可接入 Redis / local / mock。
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, data []byte, ttl int) error
	Del(ctx context.Context, key string) error
}

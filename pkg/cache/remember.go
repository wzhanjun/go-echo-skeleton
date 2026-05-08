package cache

import (
	"context"
	"encoding/json"

	"golang.org/x/sync/singleflight"
)

var sf singleflight.Group

// Remember 读穿透缓存：先查缓存 → 命中返回 → 单飞加载 → 回填缓存 → 返回。
// 同一 key 的并发请求只会执行一次 load，其余等待其结果。
func Remember[T any](ctx context.Context, c Cache, key string, ttl int, load func() (T, error)) (T, error) {
	cached, err := c.Get(ctx, key)
	if err == nil && len(cached) > 0 {
		var v T
		if e := json.Unmarshal(cached, &v); e == nil {
			return v, nil
		}
	}

	result, err, _ := sf.Do(key, func() (interface{}, error) {
		return load()
	})
	if err != nil {
		var zero T
		return zero, err
	}

	val := result.(T)
	if data, e := json.Marshal(val); e == nil {
		_ = c.Set(ctx, key, data, ttl)
	}
	return val, nil
}

// Key 构造带命名空间的缓存键，如 Key("user", "list") → "user:list"
func Key(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}
	s := parts[0]
	for _, p := range parts[1:] {
		s += ":" + p
	}
	return s
}

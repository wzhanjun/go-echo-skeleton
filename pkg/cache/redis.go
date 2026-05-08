package cache

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wzhanjun/go-echo-skeleton/pkg/config"
	slog "github.com/wzhanjun/log-service/client"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

const redisTimeout = 5 * time.Second

func GetRedisConn() *redis.Client {

	once.Do(func() {
		if len(config.Cfg.Redis.Sentinel.Nodes) >= 3 {
			redisClient = redis.NewFailoverClient(&redis.FailoverOptions{
				MasterName:       config.Cfg.Redis.Sentinel.MasterName,
				SentinelAddrs:    config.Cfg.Redis.Sentinel.Nodes,
				Password:         config.Cfg.Redis.Sentinel.Password,
				SentinelPassword: config.Cfg.Redis.Sentinel.Password,
				DB:               config.Cfg.Redis.DB,
				DialTimeout:      redisTimeout,
				ReadTimeout:      redisTimeout,
				WriteTimeout:     redisTimeout,
				ConnMaxIdleTime:  60 * time.Second,
				PoolSize:         200,
				MaxIdleConns:     50,
				MinIdleConns:     10,
			})
		} else {
			redisClient = redis.NewClient(&redis.Options{
				Addr:            config.Cfg.Redis.Addr,
				Password:        config.Cfg.Redis.Password,
				DB:              config.Cfg.Redis.DB,
				DialTimeout:     redisTimeout,
				ReadTimeout:     redisTimeout,
				WriteTimeout:    redisTimeout,
				ConnMaxIdleTime: 60 * time.Second,
				PoolSize:        200,
				MinIdleConns:    50,
			})
		}

		pingCtx, cancel := context.WithTimeout(context.Background(), redisTimeout)
		defer cancel()
		_, err := redisClient.Ping(pingCtx).Result()
		if err != nil {
			slog.Label("redis").Panicf("redis连接失败, ping err:%+v", err)
			panic(err)
		}
	})

	return redisClient
}

func GetRedisLock(key string, second int) (bool, error) {
	client := GetRedisConn()	

	ok, err := client.SetNX(context.Background(), "gameads:lock:"+key, true, time.Second*time.Duration(second)).Result()
	if err != nil {
		slog.Label("redis").Errorf("redis加锁失败, key:%s err:%+v", key, err)
		return false, err
	}

	return ok, nil
}

func RedisUnLock(key string) error {
	client := GetRedisConn()

	_, err := client.Del(context.Background(), "gameads:lock:"+key).Result()
	if err != nil {
		slog.Label("redis").Errorf("redis释放锁失败, key:%s err:%+v", key, err)
		return err
	}

	return nil
}

var incrExpire = redis.NewScript(`
local current

current = redis.call("incr", KEYS[1])
if current == 1 then
    redis.call("expire",KEYS[1],ARGV[1])
end

return current
`)

func IncrWithExpire(cxt context.Context, key string, seconds int64) (int64, error) {
	return incrExpire.Run(cxt, GetRedisConn(), []string{key}, seconds).Int64()
}

package controller

import (
	"context"
	"time"

	"server/config"

	redis "github.com/redis/go-redis/v9"
)

const redisConfigKey = "config:kv"

var rdb *redis.Client

func InitRedis() {
	if rdb != nil {
		return
	}
	cfg := config.GetConfig()
	addr := cfg.Redis.Addr
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
}

func RedisSetConfig(ctx context.Context, kv map[string]string) error {
	InitRedis()
	if len(kv) == 0 {
		return nil
	}
	// 使用Pipeline批量设置
	pipe := rdb.Pipeline()
	for k, v := range kv {
		pipe.HSet(ctx, redisConfigKey, k, v)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func RedisGetAllConfig(ctx context.Context) (map[string]string, error) {
	InitRedis()
	res, err := rdb.HGetAll(ctx, redisConfigKey).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func PingRedis(ctx context.Context) error {
	InitRedis()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return rdb.Ping(ctx).Err()
}

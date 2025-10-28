package redisx

import (
	"context"

	"server/config"

	redis "github.com/redis/go-redis/v9"
)

const ConfigHashKey = "config:kv"

var client *redis.Client

func Init() {
	if client != nil {
		return
	}
	cfg := config.GetConfig().Redis
	addr := cfg.Addr
	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	client = redis.NewClient(&redis.Options{Addr: addr, Password: cfg.Password, DB: cfg.DB})
}

func HSet(ctx context.Context, key string, kv map[string]string) error {
	Init()
	if len(kv) == 0 {
		return nil
	}
	return client.HSet(ctx, key, kv).Err()
}

func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	Init()
	return client.HGetAll(ctx, key).Result()
}

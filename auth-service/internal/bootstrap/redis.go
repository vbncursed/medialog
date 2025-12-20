package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vbncursed/medialog/auth-service/config"
)

func InitRedis(cfg *config.Config) *redis.Client {
	addr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Быстрый ping при старте (fail-open логика для auth не нужна — если Redis нужен для RL, лучше знать сразу).
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Error("redis ping failed", "addr", addr, "err", err)
		// Не паникуем — сервис может работать и без Redis (RL будет fail-open).
	}

	return rdb
}

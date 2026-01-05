package bootstrap

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/vbncursed/medialog/auth/config"
)

func InitRedis(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return rdb
}

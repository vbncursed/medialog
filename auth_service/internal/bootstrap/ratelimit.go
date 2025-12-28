package bootstrap

import (
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vbncursed/medialog/auth_service/config"
	server "github.com/vbncursed/medialog/auth_service/internal/api/auth_service_api"
)

func InitAuthRateLimiters(redisClient *redis.Client, cfg *config.Config) (server.RateLimiter, server.RateLimiter, server.RateLimiter) {
	return server.NewRedisRateLimiter(redisClient, "login", cfg.Auth.RateLimitLoginPerMinute, time.Minute),
		server.NewRedisRateLimiter(redisClient, "register", cfg.Auth.RateLimitRegisterPerMinute, time.Minute),
		server.NewRedisRateLimiter(redisClient, "refresh", cfg.Auth.RateLimitRefreshPerMinute, time.Minute)
}

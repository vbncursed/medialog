package bootstrap

import (
	"time"

	"github.com/redis/go-redis/v9"
	server "github.com/vbncursed/medialog/auth-service/internal/api/auth_service_api"
)

func InitAuthRateLimiters(redisClient *redis.Client, loginLimitPerMinute, registerLimitPerMinute int) (server.RateLimiter, server.RateLimiter) {
	return server.NewRedisRateLimiter(redisClient, "login", loginLimitPerMinute, time.Minute),
		server.NewRedisRateLimiter(redisClient, "register", registerLimitPerMinute, time.Minute)
}



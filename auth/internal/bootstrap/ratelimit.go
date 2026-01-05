package bootstrap

import (
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vbncursed/medialog/auth/config"
	"github.com/vbncursed/medialog/auth/internal/api/auth_service_api"
)

func InitRateLimiters(rdb *redis.Client, cfg *config.Config) (loginLimiter, registerLimiter, refreshLimiter auth_service_api.RateLimiter) {
	window := time.Minute

	loginLimiter = auth_service_api.NewRedisRateLimiter(
		rdb,
		"login",
		cfg.Auth.RateLimitLoginPerMinute,
		window,
	)

	registerLimiter = auth_service_api.NewRedisRateLimiter(
		rdb,
		"register",
		cfg.Auth.RateLimitRegisterPerMinute,
		window,
	)

	refreshLimiter = auth_service_api.NewRedisRateLimiter(
		rdb,
		"refresh",
		cfg.Auth.RateLimitRefreshPerMinute,
		window,
	)

	return loginLimiter, registerLimiter, refreshLimiter
}


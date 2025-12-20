package bootstrap

import (
	"github.com/redis/go-redis/v9"
	"github.com/vbncursed/medialog/auth-service/config"
	server "github.com/vbncursed/medialog/auth-service/internal/api/auth_service_api"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService"
)

func InitAuthServiceAPI(authService *authService.AuthService, redisClient *redis.Client, cfg *config.Config) *server.AuthServiceAPI {
	return server.NewAuthServiceAPI(authService, redisClient, cfg.Auth.RateLimitLoginPerMinute, cfg.Auth.RateLimitRegisterPerMinute)
}

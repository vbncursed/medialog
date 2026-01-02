package bootstrap

import (
	server "github.com/vbncursed/medialog/auth_service/internal/api/auth_service_api"
	"github.com/vbncursed/medialog/auth_service/config"
	"github.com/vbncursed/medialog/auth_service/internal/services/auth_service"
)

func InitAuthServiceAPI(authService *auth_service.AuthService, cfg *config.Config, loginLimiter, registerLimiter, refreshLimiter server.RateLimiter) *server.AuthServiceAPI {
	return server.NewAuthServiceAPI(authService, cfg.Auth.JWTSecret, loginLimiter, registerLimiter, refreshLimiter)
}

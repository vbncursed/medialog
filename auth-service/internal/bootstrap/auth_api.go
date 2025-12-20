package bootstrap

import (
	"github.com/vbncursed/medialog/auth-service/config"
	server "github.com/vbncursed/medialog/auth-service/internal/api/auth_service_api"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService"
)

func InitAuthServiceAPI(authService *authService.AuthService, cfg *config.Config) *server.AuthServiceAPI {
	return server.NewAuthServiceAPI(authService, cfg.Auth.RateLimitLoginPerMinute, cfg.Auth.RateLimitRegisterPerMinute)
}

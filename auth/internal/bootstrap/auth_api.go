package bootstrap

import (
	"github.com/vbncursed/medialog/auth/config"
	"github.com/vbncursed/medialog/auth/internal/api/auth_service_api"
	"github.com/vbncursed/medialog/auth/internal/services/auth_service"
)

func InitAuthServiceAPI(
	authService *auth_service.AuthService,
	cfg *config.Config,
	loginLimiter, registerLimiter, refreshLimiter auth_service_api.RateLimiter,
) *auth_service_api.AuthServiceAPI {
	return auth_service_api.NewAuthServiceAPI(
		authService,
		cfg.Auth.JWTSecret,
		loginLimiter,
		registerLimiter,
		refreshLimiter,
	)
}


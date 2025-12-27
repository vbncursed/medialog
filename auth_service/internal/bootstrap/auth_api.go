package bootstrap

import (
	server "github.com/vbncursed/medialog/auth-service/internal/api/auth_service_api"
	"github.com/vbncursed/medialog/auth-service/internal/services/auth_service"
)

func InitAuthServiceAPI(authService *auth_service.AuthService, loginLimiter, registerLimiter server.RateLimiter) *server.AuthServiceAPI {
	return server.NewAuthServiceAPI(authService, loginLimiter, registerLimiter)
}

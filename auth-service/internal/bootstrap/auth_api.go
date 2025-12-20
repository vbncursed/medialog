package bootstrap

import (
	server "github.com/vbncursed/medialog/auth-service/internal/api/auth_service_api"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService"
)

func InitAuthServiceAPI(authService *authService.AuthService, loginLimiter, registerLimiter server.RateLimiter) *server.AuthServiceAPI {
	return server.NewAuthServiceAPI(authService, loginLimiter, registerLimiter)
}

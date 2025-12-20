package bootstrap

import (
	"github.com/vbncursed/medialog/auth-service/config"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService"
	"github.com/vbncursed/medialog/auth-service/internal/storage/pgstorage"
)

func InitAuthService(storage *pgstorage.PGstorage, cfg *config.Config) *authService.AuthService {
	return authService.NewAuthService(storage, cfg.Auth.JWTSecret, cfg.Auth.AccessTTLSeconds, cfg.Auth.RefreshTTLSeconds)
}

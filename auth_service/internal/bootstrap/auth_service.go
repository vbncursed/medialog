package bootstrap

import (
	"github.com/vbncursed/medialog/auth-service/config"
	"github.com/vbncursed/medialog/auth-service/internal/services/auth_service"
	"github.com/vbncursed/medialog/auth-service/internal/storage/auth_storage"
)

func InitAuthService(storage *auth_storage.AuthStorage, cfg *config.Config) *auth_service.AuthService {
	return auth_service.NewAuthService(storage, cfg.Auth.JWTSecret, cfg.Auth.AccessTTLSeconds, cfg.Auth.RefreshTTLSeconds)
}

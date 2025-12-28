package bootstrap

import (
	"github.com/vbncursed/medialog/auth_service/config"
	"github.com/vbncursed/medialog/auth_service/internal/services/auth_service"
	"github.com/vbncursed/medialog/auth_service/internal/storage/auth_storage"
	"github.com/vbncursed/medialog/auth_service/internal/storage/session_storage"
)

func InitAuthService(authStorage *auth_storage.AuthStorage, sessionStorage *session_storage.SessionStorage, cfg *config.Config) *auth_service.AuthService {
	return auth_service.NewAuthService(authStorage, sessionStorage, cfg.Auth.JWTSecret, cfg.Auth.AccessTTLSeconds, cfg.Auth.RefreshTTLSeconds)
}

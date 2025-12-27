package auth_service_api

import (
	"context"

	"github.com/vbncursed/medialog/auth-service/internal/models"
	"github.com/vbncursed/medialog/auth-service/internal/pb/auth_api"
	"github.com/vbncursed/medialog/auth-service/internal/services/auth_service"
)

type authService interface {
	Register(ctx context.Context, in models.RegisterInput) (*auth_service.AuthInfo, error)
	Login(ctx context.Context, in models.LoginInput) (*auth_service.AuthInfo, error)
	Refresh(ctx context.Context, in models.RefreshInput) (*auth_service.AuthInfo, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, refreshToken string) error
}

// AuthServiceAPI реализует grpc AuthServiceServer.
type AuthServiceAPI struct {
	auth_api.UnimplementedAuthServiceServer
	authService     authService
	loginLimiter    RateLimiter
	registerLimiter RateLimiter
}

type denyAllLimiter struct{}

func (denyAllLimiter) Allow(ctx context.Context, key string) bool { return false }

func NewAuthServiceAPI(authService *auth_service.AuthService, loginLimiter, registerLimiter RateLimiter) *AuthServiceAPI {
	if loginLimiter == nil {
		loginLimiter = denyAllLimiter{}
	}
	if registerLimiter == nil {
		registerLimiter = denyAllLimiter{}
	}
	return &AuthServiceAPI{
		authService:     authService,
		loginLimiter:    loginLimiter,
		registerLimiter: registerLimiter,
	}
}

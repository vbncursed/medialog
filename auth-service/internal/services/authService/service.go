package authService

import (
	"context"

	"github.com/vbncursed/medialog/auth-service/internal/models"
)

// Service — интерфейс сервисного слоя auth-service (удобен для моков/mockery).
type Service interface {
	Register(ctx context.Context, in models.RegisterInput) (*AuthResult, error)
	Login(ctx context.Context, in models.LoginInput) (*AuthResult, error)
	Refresh(ctx context.Context, in models.RefreshInput) (*AuthResult, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, refreshToken string) error
}

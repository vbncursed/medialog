package authService

import "context"

// Service — интерфейс сервисного слоя auth-service (удобен для моков/mockery).
type Service interface {
	Register(ctx context.Context, email, password, userAgent, ip string) (*AuthResult, error)
	Login(ctx context.Context, email, password, userAgent, ip string) (*AuthResult, error)
	Refresh(ctx context.Context, refreshToken, userAgent, ip string) (*AuthResult, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, refreshToken string) error
}

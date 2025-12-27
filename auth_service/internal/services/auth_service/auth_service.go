package auth_service

import (
	"context"
	"time"

	"github.com/vbncursed/medialog/auth-service/internal/models"
)

type AuthStorage interface {
	CreateUser(ctx context.Context, email string, passwordHash string) (uint64, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	CreateSession(ctx context.Context, userID uint64, refreshHash []byte, expiresAt time.Time, userAgent, ip string) (uint64, error)
	GetSessionByRefreshHash(ctx context.Context, refreshHash []byte) (*models.Session, error)
	RevokeSessionByID(ctx context.Context, sessionID uint64, revokedAt time.Time) error
	RevokeAllSessionsByUserID(ctx context.Context, userID uint64, revokedAt time.Time) error
}

type Service interface {
	Register(ctx context.Context, in models.RegisterInput) (*AuthInfo, error)
	Login(ctx context.Context, in models.LoginInput) (*AuthInfo, error)
	Refresh(ctx context.Context, in models.RefreshInput) (*AuthInfo, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, refreshToken string) error
}

type AuthService struct {
	authStorage AuthStorage

	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type AuthInfo struct {
	UserID       uint64
	AccessToken  string
	RefreshToken string
}

var (
	tokenToHashFn     = tokenToHash
	newAccessTokenFn  = newAccessToken
	newRefreshTokenFn = newRefreshToken
)

func NewAuthService(authStorage AuthStorage, jwtSecret string, accessTTLSeconds, refreshTTLSeconds int64) *AuthService {
	return &AuthService{
		authStorage: authStorage,
		jwtSecret:   jwtSecret,
		accessTTL:   time.Duration(accessTTLSeconds) * time.Second,
		refreshTTL:  time.Duration(refreshTTLSeconds) * time.Second,
	}
}

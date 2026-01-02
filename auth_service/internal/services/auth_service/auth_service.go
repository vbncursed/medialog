package auth_service

import (
	"context"
	"time"

	"github.com/vbncursed/medialog/auth_service/internal/models"
)

type AuthStorage interface {
	CreateUser(ctx context.Context, email string, passwordHash string) (uint64, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, userID uint64) (*models.User, error)
	UpdateUserRole(ctx context.Context, userID uint64, role string) error
}

type SessionStorage interface {
	CreateSession(ctx context.Context, userID uint64, refreshHash []byte, expiresAt time.Time, userAgent, ip string) error
	GetSessionByRefreshHash(ctx context.Context, refreshHash []byte) (*models.Session, error)
	RevokeSessionByRefreshHash(ctx context.Context, refreshHash []byte) error
	RevokeAllSessionsByUserID(ctx context.Context, userID uint64) error
}

type AuthService struct {
	authStorage    AuthStorage
	sessionStorage SessionStorage

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

func NewAuthService(authStorage AuthStorage, sessionStorage SessionStorage, jwtSecret string, accessTTLSeconds, refreshTTLSeconds int64) *AuthService {
	return &AuthService{
		authStorage:    authStorage,
		sessionStorage: sessionStorage,
		jwtSecret:      jwtSecret,
		accessTTL:      time.Duration(accessTTLSeconds) * time.Second,
		refreshTTL:     time.Duration(refreshTTLSeconds) * time.Second,
	}
}

package authService

import "time"

// AuthService — сервисный слой auth.
// В этом файле держим “каркас”: типы, конструктор и test hooks.
// Реальные методы разнесены по отдельным файлам по use-case’ам.
type AuthService struct {
	storage Storage

	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type AuthInfo struct {
	UserID       uint64
	AccessToken  string
	RefreshToken string
}

// test hooks for hard-to-trigger branches
var (
	tokenToHashFn     = tokenToHash
	newAccessTokenFn  = newAccessToken
	newRefreshTokenFn = newRefreshToken
)

func NewAuthService(storage Storage, jwtSecret string, accessTTLSeconds, refreshTTLSeconds int64) *AuthService {
	return &AuthService{
		storage:    storage,
		jwtSecret:  jwtSecret,
		accessTTL:  time.Duration(accessTTLSeconds) * time.Second,
		refreshTTL: time.Duration(refreshTTLSeconds) * time.Second,
	}
}

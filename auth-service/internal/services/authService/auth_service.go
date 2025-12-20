package authService

import "time"

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

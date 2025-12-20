package authService_test

import (
	"crypto/sha256"

	"github.com/vbncursed/medialog/auth-service/internal/models"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService"
)

const (
	testJWTSecret  = "secret"
	testAccessTTL  = int64(60)
	testRefreshTTL = int64(3600)
)

func sha256b(v string) []byte {
	sum := sha256.Sum256([]byte(v))
	return sum[:]
}

func newTestService(st authService.Storage) *authService.AuthService {
	return authService.NewAuthService(st, testJWTSecret, testAccessTTL, testRefreshTTL)
}

type emailPasswordIPInput interface {
	~struct {
		Email     string
		Password  string
		UserAgent string
		IP        string
	}
}

func inEmailPass[T emailPasswordIPInput](email, password string) T {
	return T{
		Email:    email,
		Password: password,
		IP:       "127.0.0.1",
	}
}

func inEmailPassWithUA[T emailPasswordIPInput](email, password, userAgent string) T {
	return T{
		Email:     email,
		Password:  password,
		UserAgent: userAgent,
		IP:        "127.0.0.1",
	}
}

func refreshIn(refreshToken string) models.RefreshInput {
	return models.RefreshInput{
		RefreshToken: refreshToken,
		UserAgent:    "ua",
		IP:           "ip",
	}
}

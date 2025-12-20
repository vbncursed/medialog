package authService_test

import (
	"context"
	"crypto/sha256"
	"testing"

	"github.com/vbncursed/medialog/auth-service/internal/models"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService/mocks"
)

const (
	testJWTSecret  = "secret"
	testAccessTTL  = int64(60)
	testRefreshTTL = int64(3600)
)

func bg() context.Context { return context.Background() }

func sha256b(v string) []byte {
	sum := sha256.Sum256([]byte(v))
	return sum[:]
}

func newTestService(st authService.Storage) *authService.AuthService {
	return authService.NewAuthService(st, testJWTSecret, testAccessTTL, testRefreshTTL)
}

// setup — общий boilerplate для unit-тестов сервиса.
// Возвращает готовый сервис и мок хранилища, чтобы в тесте осталось только EXPECT + вызов + assert.
func setup(t *testing.T) (*authService.AuthService, *mocks.Storage) {
	t.Helper()
	st := mocks.NewStorage(t)
	return newTestService(st), st
}

type emailPasswordIPInput interface {
	~struct {
		Email     string
		Password  string
		UserAgent string
		IP        string
	}
}

// inEmailPass — общий конструктор для LoginInput/RegisterInput в тестах.
// Дефолты: IP=127.0.0.1, UserAgent="".
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

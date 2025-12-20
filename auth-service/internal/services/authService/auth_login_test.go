package authService

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vbncursed/medialog/auth-service/internal/models"
)

func TestAuthService_Login_InvalidArgs(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	_, err := svc.Login(context.Background(), models.LoginInput{
		Email:    "bad",
		Password: "short",
		IP:       "127.0.0.1",
	})
	require.ErrorIs(t, err, ErrInvalidArgument)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	_, err := svc.Login(context.Background(), models.LoginInput{
		Email:    "a@b.com",
		Password: "Password123",
		IP:       "127.0.0.1",
	})
	require.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	st := newFakeStorage()
	passHash, err := passwordHash("Password123")
	require.NoError(t, err)
	_, _ = st.CreateUser(context.Background(), "a@b.com", passHash)

	svc := NewAuthService(st, "secret", 60, 3600)
	_, err = svc.Login(context.Background(), models.LoginInput{
		Email:    "a@b.com",
		Password: "Password124",
		IP:       "127.0.0.1",
	})
	require.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestAuthService_Login_Success(t *testing.T) {
	st := newFakeStorage()
	passHash, err := passwordHash("Password123")
	require.NoError(t, err)
	_, _ = st.CreateUser(context.Background(), "a@b.com", passHash)

	svc := NewAuthService(st, "secret", 60, 3600)
	res, err := svc.Login(context.Background(), models.LoginInput{
		Email:     "a@b.com",
		Password:  "Password123",
		UserAgent: "ua",
		IP:        "127.0.0.1",
	})
	require.NoError(t, err)
	require.NotEmpty(t, res.AccessToken)
	require.NotEmpty(t, res.RefreshToken)
}

func TestAuthService_Login_StorageError(t *testing.T) {
	st := newFakeStorage()
	st.errGetUser = errors.New("db fail")
	svc := NewAuthService(st, "secret", 60, 3600)

	_, err := svc.Login(context.Background(), models.LoginInput{
		Email:     "a@b.com",
		Password:  "Password123",
		UserAgent: "ua",
		IP:        "127.0.0.1",
	})
	require.Error(t, err)
}

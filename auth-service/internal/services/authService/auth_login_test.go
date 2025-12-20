package authService

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthService_Login_InvalidArgs(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	_, err := svc.Login(context.Background(), "bad", "short", "", "127.0.0.1")
	require.ErrorIs(t, err, ErrInvalidArgument)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	_, err := svc.Login(context.Background(), "a@b.com", "Password123", "", "127.0.0.1")
	require.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	st := newFakeStorage()
	passHash, err := hashPassword("Password123")
	require.NoError(t, err)
	_, _ = st.CreateUser(context.Background(), "a@b.com", passHash)

	svc := NewAuthService(st, "secret", 60, 3600)
	_, err = svc.Login(context.Background(), "a@b.com", "Password124", "", "127.0.0.1")
	require.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestAuthService_Login_Success(t *testing.T) {
	st := newFakeStorage()
	passHash, err := hashPassword("Password123")
	require.NoError(t, err)
	_, _ = st.CreateUser(context.Background(), "a@b.com", passHash)

	svc := NewAuthService(st, "secret", 60, 3600)
	res, err := svc.Login(context.Background(), "a@b.com", "Password123", "ua", "127.0.0.1")
	require.NoError(t, err)
	require.NotEmpty(t, res.AccessToken)
	require.NotEmpty(t, res.RefreshToken)
}

func TestAuthService_Login_StorageError(t *testing.T) {
	st := newFakeStorage()
	st.errGetUser = errors.New("db fail")
	svc := NewAuthService(st, "secret", 60, 3600)

	_, err := svc.Login(context.Background(), "a@b.com", "Password123", "ua", "127.0.0.1")
	require.Error(t, err)
}

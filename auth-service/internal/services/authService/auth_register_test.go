package authService

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthService_Register_InvalidArgs(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	_, err := svc.Register(context.Background(), "bad", "short", "", "127.0.0.1")
	require.ErrorIs(t, err, ErrInvalidArgument)
}

func TestAuthService_Register_PasswordComplexity(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	cases := []string{
		"password123", // no upper
		"PASSWORD123", // no lower
		"Password",    // no digit
		"Passw1",      // too short
	}

	for _, pwd := range cases {
		_, err := svc.Register(context.Background(), "a@b.com", pwd, "", "127.0.0.1")
		require.ErrorIs(t, err, ErrInvalidArgument, "pwd=%q", pwd)
	}
}

func TestAuthService_Register_EmailExists(t *testing.T) {
	st := newFakeStorage()
	_, _ = st.CreateUser(context.Background(), "a@b.com", "hash")

	svc := NewAuthService(st, "secret", 60, 3600)
	_, err := svc.Register(context.Background(), "a@b.com", "Password123", "", "127.0.0.1")
	require.ErrorIs(t, err, ErrEmailAlreadyExists)
}

func TestAuthService_Register_StorageLookupError(t *testing.T) {
	st := newFakeStorage()
	st.errGetUser = errors.New("boom")

	svc := NewAuthService(st, "secret", 60, 3600)
	_, err := svc.Register(context.Background(), "a@b.com", "Password123", "", "127.0.0.1")
	require.Error(t, err)
}

func TestAuthService_Register_CreateUserErrorMappedToAlreadyExists(t *testing.T) {
	st := newFakeStorage()
	st.errCreateUser = errors.New("db down")

	svc := NewAuthService(st, "secret", 60, 3600)
	_, err := svc.Register(context.Background(), "a@b.com", "Password123", "", "127.0.0.1")
	require.ErrorIs(t, err, ErrEmailAlreadyExists)
}

func TestAuthService_Register_Success(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	res, err := svc.Register(context.Background(), "a@b.com", "Password123", "ua", "127.0.0.1")
	require.NoError(t, err)
	require.NotZero(t, res.UserID)
	require.NotEmpty(t, res.AccessToken)
	require.NotEmpty(t, res.RefreshToken)
}

func TestAuthService_Register_HashPasswordError(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	old := bcryptGenerate
	bcryptGenerate = func(_ []byte, _ int) ([]byte, error) { return nil, errors.New("bcrypt fail") }
	t.Cleanup(func() { bcryptGenerate = old })

	_, err := svc.Register(context.Background(), "a@b.com", "Password123", "", "127.0.0.1")
	require.Error(t, err)
}

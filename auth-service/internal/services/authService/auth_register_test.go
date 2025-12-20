package authService_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/vbncursed/medialog/auth-service/internal/models"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService"
	pguserstorage "github.com/vbncursed/medialog/auth-service/internal/storage/pgUserStorage"
	"gotest.tools/v3/assert"
)

func TestAuthService_Register_InvalidArgs(t *testing.T) {
	svc, _ := setup(t)
	_, gotErr := svc.Register(bg(), inEmailPass[models.RegisterInput]("bad", "short"))
	wantErr := authService.ErrInvalidArgument
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_Register_PasswordComplexity(t *testing.T) {
	svc, _ := setup(t)

	cases := []string{
		"password123", // no upper
		"PASSWORD123", // no lower
		"Password",    // no digit
		"Passw1",      // too short
	}

	for _, pwd := range cases {
		_, gotErr := svc.Register(bg(), inEmailPass[models.RegisterInput]("a@b.com", pwd))
		wantErr := authService.ErrInvalidArgument
		assert.ErrorIs(t, gotErr, wantErr)
	}
}

func TestAuthService_Register_EmailExists(t *testing.T) {
	svc, st := setup(t)
	st.EXPECT().
		GetUserByEmail(bg(), "a@b.com").
		Return(&models.User{ID: 1, Email: "a@b.com", PasswordHash: "hash"}, nil)
	_, gotErr := svc.Register(bg(), inEmailPass[models.RegisterInput]("a@b.com", "Password123"))
	wantErr := authService.ErrEmailAlreadyExists
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_Register_StorageLookupError(t *testing.T) {
	wantErr := errors.New("boom")
	svc, st := setup(t)
	st.EXPECT().GetUserByEmail(bg(), "a@b.com").Return(nil, wantErr)

	_, gotErr := svc.Register(bg(), inEmailPass[models.RegisterInput]("a@b.com", "Password123"))
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_Register_CreateUserErrorMappedToAlreadyExists(t *testing.T) {
	svc, st := setup(t)
	st.EXPECT().GetUserByEmail(bg(), "a@b.com").Return(nil, pguserstorage.ErrUserNotFound)
	st.EXPECT().CreateUser(mock.Anything, "a@b.com", mock.Anything).Return(uint64(0), errors.New("db down"))

	_, gotErr := svc.Register(bg(), inEmailPass[models.RegisterInput]("a@b.com", "Password123"))
	wantErr := authService.ErrEmailAlreadyExists
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_Register_Success(t *testing.T) {
	svc, st := setup(t)
	st.EXPECT().GetUserByEmail(bg(), "a@b.com").Return(nil, pguserstorage.ErrUserNotFound)
	st.EXPECT().CreateUser(mock.Anything, "a@b.com", mock.Anything).Return(uint64(1), nil)
	st.EXPECT().CreateSession(mock.Anything, uint64(1), mock.Anything, mock.Anything, "ua", "127.0.0.1").Return(uint64(1), nil)

	got, gotErr := svc.Register(bg(), inEmailPassWithUA[models.RegisterInput]("a@b.com", "Password123", "ua"))
	assert.NilError(t, gotErr)

	wantUserIDNonZero := true
	wantAccessNonEmpty := true
	wantRefreshNonEmpty := true
	assert.Equal(t, got.UserID != 0, wantUserIDNonZero)
	assert.Equal(t, got.AccessToken != "", wantAccessNonEmpty)
	assert.Equal(t, got.RefreshToken != "", wantRefreshNonEmpty)
}

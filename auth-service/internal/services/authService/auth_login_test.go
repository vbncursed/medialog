package authService_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/vbncursed/medialog/auth-service/internal/models"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService"
	pguserstorage "github.com/vbncursed/medialog/auth-service/internal/storage/pgUserStorage"
	"golang.org/x/crypto/bcrypt"
	"gotest.tools/v3/assert"
)

func TestAuthService_Login_InvalidArgs(t *testing.T) {
	svc, _ := setup(t)
	_, gotErr := svc.Login(bg(), inEmailPass[models.LoginInput]("bad", "short"))
	wantErr := authService.ErrInvalidArgument
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	svc, st := setup(t)
	st.EXPECT().GetUserByEmail(bg(), "a@b.com").Return(nil, pguserstorage.ErrUserNotFound)
	_, gotErr := svc.Login(bg(), inEmailPass[models.LoginInput]("a@b.com", "Password123"))
	wantErr := authService.ErrInvalidCredentials
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	passHashBytes, err := bcrypt.GenerateFromPassword([]byte("Password123"), bcrypt.DefaultCost)
	assert.NilError(t, err)
	passHash := string(passHashBytes)

	svc, st := setup(t)
	st.EXPECT().
		GetUserByEmail(bg(), "a@b.com").
		Return(&models.User{ID: 1, Email: "a@b.com", PasswordHash: passHash}, nil)

	_, gotErr := svc.Login(bg(), inEmailPass[models.LoginInput]("a@b.com", "Password124"))
	wantErr := authService.ErrInvalidCredentials
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_Login_Success(t *testing.T) {
	passHashBytes, err := bcrypt.GenerateFromPassword([]byte("Password123"), bcrypt.DefaultCost)
	assert.NilError(t, err)
	passHash := string(passHashBytes)

	svc, st := setup(t)
	st.EXPECT().
		GetUserByEmail(bg(), "a@b.com").
		Return(&models.User{ID: 1, Email: "a@b.com", PasswordHash: passHash}, nil)
	st.EXPECT().CreateSession(mock.Anything, uint64(1), mock.Anything, mock.Anything, "ua", "127.0.0.1").Return(uint64(1), nil)

	got, gotErr := svc.Login(bg(), models.LoginInput{
		Email:     "a@b.com",
		Password:  "Password123",
		UserAgent: "ua",
		IP:        "127.0.0.1",
	})
	assert.NilError(t, gotErr)

	wantAccessNonEmpty := true
	wantRefreshNonEmpty := true
	gotAccessNonEmpty := got.AccessToken != ""
	gotRefreshNonEmpty := got.RefreshToken != ""
	assert.Equal(t, gotAccessNonEmpty, wantAccessNonEmpty)
	assert.Equal(t, gotRefreshNonEmpty, wantRefreshNonEmpty)
}

func TestAuthService_Login_StorageError(t *testing.T) {
	wantErr := errors.New("db fail")
	svc, st := setup(t)
	st.EXPECT().GetUserByEmail(bg(), "a@b.com").Return(nil, wantErr)
	_, gotErr := svc.Login(bg(), models.LoginInput{
		Email:     "a@b.com",
		Password:  "Password123",
		UserAgent: "ua",
		IP:        "127.0.0.1",
	})
	assert.ErrorIs(t, gotErr, wantErr)
}

package authService_test

import (
	"errors"

	"github.com/stretchr/testify/mock"
	"github.com/vbncursed/medialog/auth-service/internal/models"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService"
	pguserstorage "github.com/vbncursed/medialog/auth-service/internal/storage/pgUserStorage"
	"golang.org/x/crypto/bcrypt"
	"gotest.tools/v3/assert"
)

func (s *AuthServiceSuite) TestLogin_InvalidArgs() {
	_, gotErr := s.svc.Login(s.ctx, inEmailPass[models.LoginInput]("bad", "short"))
	wantErr := authService.ErrInvalidArgument
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestLogin_UserNotFound() {
	s.st.EXPECT().GetUserByEmail(s.ctx, "a@b.com").Return(nil, pguserstorage.ErrUserNotFound)
	_, gotErr := s.svc.Login(s.ctx, inEmailPass[models.LoginInput]("a@b.com", "Password123"))
	wantErr := authService.ErrInvalidCredentials
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestLogin_WrongPassword() {
	passHashBytes, err := bcrypt.GenerateFromPassword([]byte("Password123"), bcrypt.DefaultCost)
	assert.NilError(s.T(), err)
	passHash := string(passHashBytes)

	s.st.EXPECT().
		GetUserByEmail(s.ctx, "a@b.com").
		Return(&models.User{ID: 1, Email: "a@b.com", PasswordHash: passHash}, nil)

	_, gotErr := s.svc.Login(s.ctx, inEmailPass[models.LoginInput]("a@b.com", "Password124"))
	wantErr := authService.ErrInvalidCredentials
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestLogin_Success() {
	passHashBytes, err := bcrypt.GenerateFromPassword([]byte("Password123"), bcrypt.DefaultCost)
	assert.NilError(s.T(), err)
	passHash := string(passHashBytes)

	s.st.EXPECT().
		GetUserByEmail(s.ctx, "a@b.com").
		Return(&models.User{ID: 1, Email: "a@b.com", PasswordHash: passHash}, nil)
	s.st.EXPECT().CreateSession(mock.Anything, uint64(1), mock.Anything, mock.Anything, "ua", "127.0.0.1").Return(uint64(1), nil)

	got, gotErr := s.svc.Login(s.ctx, models.LoginInput{
		Email:     "a@b.com",
		Password:  "Password123",
		UserAgent: "ua",
		IP:        "127.0.0.1",
	})
	assert.NilError(s.T(), gotErr)

	wantAccessNonEmpty := true
	wantRefreshNonEmpty := true
	assert.Equal(s.T(), got.AccessToken != "", wantAccessNonEmpty)
	assert.Equal(s.T(), got.RefreshToken != "", wantRefreshNonEmpty)
}

func (s *AuthServiceSuite) TestLogin_StorageError() {
	wantErr := errors.New("db fail")
	s.st.EXPECT().GetUserByEmail(s.ctx, "a@b.com").Return(nil, wantErr)
	_, gotErr := s.svc.Login(s.ctx, models.LoginInput{
		Email:     "a@b.com",
		Password:  "Password123",
		UserAgent: "ua",
		IP:        "127.0.0.1",
	})
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

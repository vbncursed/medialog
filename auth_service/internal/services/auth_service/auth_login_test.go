package auth_service_test

import (
	"errors"

	"github.com/stretchr/testify/mock"
	"github.com/vbncursed/medialog/auth_service/internal/models"
	"github.com/vbncursed/medialog/auth_service/internal/services/auth_service"
	"github.com/vbncursed/medialog/auth_service/internal/storage/auth_storage"
	"golang.org/x/crypto/bcrypt"
	"gotest.tools/v3/assert"
)

func (s *AuthServiceSuite) TestLogin_InvalidArgs() {
	_, got := s.svc.Login(s.ctx, inEmailPass[models.LoginInput]("bad", "short"))
	want := auth_service.ErrInvalidArgument
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestLogin_UserNotFound() {
	s.st.EXPECT().GetUserByEmail(s.ctx, "a@b.com").Return(nil, auth_storage.ErrUserNotFound)
	_, got := s.svc.Login(s.ctx, inEmailPass[models.LoginInput]("a@b.com", "Password123"))
	want := auth_service.ErrInvalidCredentials
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestLogin_WrongPassword() {
	passHashBytes, err := bcrypt.GenerateFromPassword([]byte("Password123"), bcrypt.DefaultCost)
	assert.NilError(s.T(), err)
	passHash := string(passHashBytes)

	s.st.EXPECT().
		GetUserByEmail(s.ctx, "a@b.com").
		Return(&models.User{ID: 1, Email: "a@b.com", PasswordHash: passHash}, nil)

	_, got := s.svc.Login(s.ctx, inEmailPass[models.LoginInput]("a@b.com", "Password124"))
	want := auth_service.ErrInvalidCredentials
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestLogin_Success() {
	passHashBytes, err := bcrypt.GenerateFromPassword([]byte("Password123"), bcrypt.DefaultCost)
	assert.NilError(s.T(), err)
	passHash := string(passHashBytes)

	s.st.EXPECT().
		GetUserByEmail(s.ctx, "a@b.com").
		Return(&models.User{ID: 1, Email: "a@b.com", PasswordHash: passHash, Role: models.RoleUser}, nil)
	s.st.EXPECT().GetUserByID(mock.Anything, uint64(1)).Return(&models.User{ID: 1, Email: "a@b.com", Role: models.RoleUser}, nil)
	s.sessSt.EXPECT().CreateSession(mock.Anything, uint64(1), mock.Anything, mock.Anything, "ua", "127.0.0.1").Return(nil)

	got, err := s.svc.Login(s.ctx, models.LoginInput{
		Email:     "a@b.com",
		Password:  "Password123",
		UserAgent: "ua",
		IP:        "127.0.0.1",
	})
	assert.NilError(s.T(), err)

	assert.Assert(s.T(), got.AccessToken != "")
	assert.Assert(s.T(), got.RefreshToken != "")
}

func (s *AuthServiceSuite) TestLogin_StorageError() {
	want := errors.New("db fail")
	s.st.EXPECT().GetUserByEmail(s.ctx, "a@b.com").Return(nil, want)
	_, got := s.svc.Login(s.ctx, models.LoginInput{
		Email:     "a@b.com",
		Password:  "Password123",
		UserAgent: "ua",
		IP:        "127.0.0.1",
	})
	assert.ErrorIs(s.T(), got, want)
}

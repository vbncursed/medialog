package authService_test

import (
	"errors"

	"github.com/stretchr/testify/mock"
	"github.com/vbncursed/medialog/auth-service/internal/models"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService"
	pguserstorage "github.com/vbncursed/medialog/auth-service/internal/storage/pgUserStorage"
	"gotest.tools/v3/assert"
)

func (s *AuthServiceSuite) TestRegister_InvalidArgs() {
	_, gotErr := s.svc.Register(s.ctx, inEmailPass[models.RegisterInput]("bad", "short"))
	wantErr := authService.ErrInvalidArgument
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestRegister_PasswordComplexity() {
	cases := []string{
		"password123", // no upper
		"PASSWORD123", // no lower
		"Password",    // no digit
		"Passw1",      // too short
	}

	for _, pwd := range cases {
		_, gotErr := s.svc.Register(s.ctx, inEmailPass[models.RegisterInput]("a@b.com", pwd))
		wantErr := authService.ErrInvalidArgument
		assert.ErrorIs(s.T(), gotErr, wantErr)
	}
}

func (s *AuthServiceSuite) TestRegister_EmailExists() {
	s.st.EXPECT().
		GetUserByEmail(s.ctx, "a@b.com").
		Return(&models.User{ID: 1, Email: "a@b.com", PasswordHash: "hash"}, nil)
	_, gotErr := s.svc.Register(s.ctx, inEmailPass[models.RegisterInput]("a@b.com", "Password123"))
	wantErr := authService.ErrEmailAlreadyExists
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestRegister_StorageLookupError() {
	wantErr := errors.New("boom")
	s.st.EXPECT().GetUserByEmail(s.ctx, "a@b.com").Return(nil, wantErr)

	_, gotErr := s.svc.Register(s.ctx, inEmailPass[models.RegisterInput]("a@b.com", "Password123"))
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestRegister_CreateUserErrorMappedToAlreadyExists() {
	s.st.EXPECT().GetUserByEmail(s.ctx, "a@b.com").Return(nil, pguserstorage.ErrUserNotFound)
	s.st.EXPECT().CreateUser(mock.Anything, "a@b.com", mock.Anything).Return(uint64(0), errors.New("db down"))

	_, gotErr := s.svc.Register(s.ctx, inEmailPass[models.RegisterInput]("a@b.com", "Password123"))
	wantErr := authService.ErrEmailAlreadyExists
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestRegister_Success() {
	s.st.EXPECT().GetUserByEmail(s.ctx, "a@b.com").Return(nil, pguserstorage.ErrUserNotFound)
	s.st.EXPECT().CreateUser(mock.Anything, "a@b.com", mock.Anything).Return(uint64(1), nil)
	s.st.EXPECT().CreateSession(mock.Anything, uint64(1), mock.Anything, mock.Anything, "ua", "127.0.0.1").Return(uint64(1), nil)

	got, gotErr := s.svc.Register(s.ctx, inEmailPassWithUA[models.RegisterInput]("a@b.com", "Password123", "ua"))
	assert.NilError(s.T(), gotErr)

	wantUserIDNonZero := true
	wantAccessNonEmpty := true
	wantRefreshNonEmpty := true
	assert.Equal(s.T(), got.UserID != 0, wantUserIDNonZero)
	assert.Equal(s.T(), got.AccessToken != "", wantAccessNonEmpty)
	assert.Equal(s.T(), got.RefreshToken != "", wantRefreshNonEmpty)
}

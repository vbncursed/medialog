package auth_service_test

import (
	"errors"

	"github.com/stretchr/testify/mock"
	"github.com/vbncursed/medialog/auth_service/internal/models"
	"github.com/vbncursed/medialog/auth_service/internal/services/auth_service"
	"github.com/vbncursed/medialog/auth_service/internal/storage/auth_storage"
	"gotest.tools/v3/assert"
)

func (s *AuthServiceSuite) TestRegister_InvalidArgs() {
	_, got := s.svc.Register(s.ctx, inEmailPass[models.RegisterInput]("bad", "short"))
	want := auth_service.ErrInvalidArgument
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestRegister_PasswordComplexity() {
	cases := []string{
		"password123", // no upper
		"PASSWORD123", // no lower
		"Password",    // no digit
		"Passw1",      // too short
	}

	for _, pwd := range cases {
		_, got := s.svc.Register(s.ctx, inEmailPass[models.RegisterInput]("a@b.com", pwd))
		want := auth_service.ErrInvalidPassword
		assert.ErrorIs(s.T(), got, want)
	}
}

func (s *AuthServiceSuite) TestRegister_EmailExists() {
	s.st.EXPECT().
		GetUserByEmail(s.ctx, "a@b.com").
		Return(&models.User{ID: 1, Email: "a@b.com", PasswordHash: "hash"}, nil)
	_, got := s.svc.Register(s.ctx, inEmailPass[models.RegisterInput]("a@b.com", "Password123"))
	want := auth_service.ErrEmailAlreadyExists
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestRegister_StorageLookupError() {
	want := errors.New("boom")
	s.st.EXPECT().GetUserByEmail(s.ctx, "a@b.com").Return(nil, want)

	_, got := s.svc.Register(s.ctx, inEmailPass[models.RegisterInput]("a@b.com", "Password123"))
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestRegister_CreateUserErrorMappedToAlreadyExists() {
	s.st.EXPECT().GetUserByEmail(s.ctx, "a@b.com").Return(nil, auth_storage.ErrUserNotFound)
	s.st.EXPECT().CreateUser(mock.Anything, "a@b.com", mock.Anything).Return(uint64(0), errors.New("db down"))

	_, got := s.svc.Register(s.ctx, inEmailPass[models.RegisterInput]("a@b.com", "Password123"))
	want := auth_service.ErrEmailAlreadyExists
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestRegister_Success() {
	s.st.EXPECT().GetUserByEmail(s.ctx, "a@b.com").Return(nil, auth_storage.ErrUserNotFound)
	s.st.EXPECT().CreateUser(mock.Anything, "a@b.com", mock.Anything).Return(uint64(1), nil)
	s.st.EXPECT().GetUserByID(mock.Anything, uint64(1)).Return(&models.User{ID: 1, Email: "a@b.com", Role: models.RoleUser}, nil)
	s.sessSt.EXPECT().CreateSession(mock.Anything, uint64(1), mock.Anything, mock.Anything, "ua", "127.0.0.1").Return(nil)

	got, err := s.svc.Register(s.ctx, inEmailPassWithUA[models.RegisterInput]("a@b.com", "Password123", "ua"))
	assert.NilError(s.T(), err)

	assert.Assert(s.T(), got.UserID != 0)
	assert.Assert(s.T(), got.AccessToken != "")
	assert.Assert(s.T(), got.RefreshToken != "")
}

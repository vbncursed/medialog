package auth_service_test

import (
	"errors"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/vbncursed/medialog/auth_service/internal/models"
	"github.com/vbncursed/medialog/auth_service/internal/services/auth_service"
	"gotest.tools/v3/assert"
)

func (s *AuthServiceSuite) TestRefresh_InvalidArgs() {
	_, got := s.svc.Refresh(s.ctx, models.RefreshInput{RefreshToken: ""})
	want := auth_service.ErrInvalidArgument
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestRefresh_SessionNotFound() {
	s.sessSt.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("tok")).Return(nil, auth_service.ErrSessionNotFound)
	_, got := s.svc.Refresh(s.ctx, refreshIn("tok"))
	want := auth_service.ErrInvalidRefreshToken
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestRefresh_GetSessionOtherError() {
	want := errors.New("db fail")
	s.sessSt.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("tok")).Return(nil, want)

	_, got := s.svc.Refresh(s.ctx, refreshIn("tok"))
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestRefresh_Revoked() {
	now := time.Now()
	rt := "refresh1"
	revokedAt := now.Add(-time.Minute)

	s.sessSt.EXPECT().
		GetSessionByRefreshHash(mock.Anything, sha256b(rt)).
		Return(&models.Session{
			ID:        10,
			UserID:    1,
			ExpiresAt: now.Add(time.Hour),
			RevokedAt: &revokedAt,
		}, nil)
	_, got := s.svc.Refresh(s.ctx, refreshIn(rt))
	want := auth_service.ErrSessionRevoked
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestRefresh_Expired() {
	now := time.Now()
	rt := "refresh1"

	s.sessSt.EXPECT().
		GetSessionByRefreshHash(mock.Anything, sha256b(rt)).
		Return(&models.Session{
			ID:        10,
			UserID:    1,
			ExpiresAt: now.Add(-time.Second),
			RevokedAt: nil,
		}, nil)
	_, got := s.svc.Refresh(s.ctx, refreshIn(rt))
	want := auth_service.ErrSessionExpired
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestRefresh_RevokeSessionError() {
	now := time.Now()
	rt := "refresh1"
	want := errors.New("revoke fail")

	s.sessSt.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b(rt)).
		Return(&models.Session{ID: 10, UserID: 1, ExpiresAt: now.Add(time.Hour)}, nil)
	s.sessSt.EXPECT().RevokeSessionByRefreshHash(mock.Anything, sha256b(rt)).Return(want)

	_, got := s.svc.Refresh(s.ctx, refreshIn(rt))
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestRefresh_Success_Rotates() {
	now := time.Now()
	rt := "refresh1"

	s.sessSt.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b(rt)).
		Return(&models.Session{ID: 10, UserID: 1, ExpiresAt: now.Add(time.Hour)}, nil)
	s.sessSt.EXPECT().RevokeSessionByRefreshHash(mock.Anything, sha256b(rt)).Return(nil)
	s.st.EXPECT().GetUserByID(mock.Anything, uint64(1)).Return(&models.User{ID: 1, Email: "a@b.com", Role: models.RoleUser}, nil)
	s.sessSt.EXPECT().CreateSession(mock.Anything, uint64(1), mock.Anything, mock.Anything, "ua", "ip").
		Return(nil)

	got, err := s.svc.Refresh(s.ctx, refreshIn(rt))
	assert.NilError(s.T(), err)

	assert.Assert(s.T(), got.AccessToken != "")
	assert.Assert(s.T(), got.RefreshToken != "")
}

func (s *AuthServiceSuite) TestLogout_InvalidArgs() {
	got := s.svc.Logout(s.ctx, "")
	want := auth_service.ErrInvalidArgument
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestLogout_SessionNotFound() {
	s.sessSt.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("tok")).Return(nil, auth_service.ErrSessionNotFound)

	got := s.svc.Logout(s.ctx, "tok")
	want := auth_service.ErrInvalidRefreshToken
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestLogout_GetSessionOtherError() {
	want := errors.New("db fail")
	s.sessSt.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("tok")).Return(nil, want)

	got := s.svc.Logout(s.ctx, "tok")
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestLogout_RevokeError() {
	want := errors.New("revoke fail")
	s.sessSt.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("refresh1")).
		Return(&models.Session{ID: 10, UserID: 1, ExpiresAt: time.Now().Add(time.Hour)}, nil)
	s.sessSt.EXPECT().RevokeSessionByRefreshHash(mock.Anything, sha256b("refresh1")).Return(want)

	got := s.svc.Logout(s.ctx, "refresh1")
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestLogout_Success() {
	s.sessSt.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("refresh1")).
		Return(&models.Session{ID: 10, UserID: 1, ExpiresAt: time.Now().Add(time.Hour)}, nil)
	s.sessSt.EXPECT().RevokeSessionByRefreshHash(mock.Anything, sha256b("refresh1")).
		Return(nil)

	got := s.svc.Logout(s.ctx, "refresh1")
	assert.NilError(s.T(), got)
}

func (s *AuthServiceSuite) TestLogoutAll_InvalidArgs() {
	got := s.svc.LogoutAll(s.ctx, "")
	want := auth_service.ErrInvalidArgument
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestLogoutAll_SessionNotFound() {
	s.sessSt.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("tok")).Return(nil, auth_service.ErrSessionNotFound)

	got := s.svc.LogoutAll(s.ctx, "tok")
	want := auth_service.ErrInvalidRefreshToken
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestLogoutAll_GetSessionOtherError() {
	want := errors.New("db fail")
	s.sessSt.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("tok")).Return(nil, want)

	got := s.svc.LogoutAll(s.ctx, "tok")
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestLogoutAll_RevokeAllError() {
	want := errors.New("revoke all fail")
	s.sessSt.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("refresh1")).
		Return(&models.Session{ID: 10, UserID: 1, ExpiresAt: time.Now().Add(time.Hour)}, nil)
	s.sessSt.EXPECT().RevokeAllSessionsByUserID(mock.Anything, uint64(1)).Return(want)

	got := s.svc.LogoutAll(s.ctx, "refresh1")
	assert.ErrorIs(s.T(), got, want)
}

func (s *AuthServiceSuite) TestLogoutAll_Success() {
	s.sessSt.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("refresh1")).
		Return(&models.Session{ID: 10, UserID: 1, ExpiresAt: time.Now().Add(time.Hour)}, nil)
	s.sessSt.EXPECT().RevokeAllSessionsByUserID(mock.Anything, uint64(1)).
		Return(nil)

	got := s.svc.LogoutAll(s.ctx, "refresh1")
	assert.NilError(s.T(), got)
}

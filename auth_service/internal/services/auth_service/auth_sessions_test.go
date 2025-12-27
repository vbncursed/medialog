package auth_service_test

import (
	"errors"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/vbncursed/medialog/auth-service/internal/models"
	"github.com/vbncursed/medialog/auth-service/internal/services/auth_service"
	"github.com/vbncursed/medialog/auth-service/internal/storage/auth_storage"
	"gotest.tools/v3/assert"
)

func (s *AuthServiceSuite) TestRefresh_InvalidArgs() {
	_, gotErr := s.svc.Refresh(s.ctx, models.RefreshInput{RefreshToken: ""})
	wantErr := auth_service.ErrInvalidArgument
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestRefresh_SessionNotFound() {
	s.st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("tok")).Return(nil, auth_storage.ErrSessionNotFound)
	_, gotErr := s.svc.Refresh(s.ctx, refreshIn("tok"))
	wantErr := auth_service.ErrInvalidRefreshToken
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestRefresh_GetSessionOtherError() {
	wantErr := errors.New("db fail")
	s.st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("tok")).Return(nil, wantErr)

	_, gotErr := s.svc.Refresh(s.ctx, refreshIn("tok"))
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestRefresh_Revoked() {
	now := time.Now()
	rt := "refresh1"
	revokedAt := now.Add(-time.Minute)

	s.st.EXPECT().
		GetSessionByRefreshHash(mock.Anything, sha256b(rt)).
		Return(&models.Session{
			ID:        10,
			UserID:    1,
			ExpiresAt: now.Add(time.Hour),
			RevokedAt: &revokedAt,
		}, nil)
	_, gotErr := s.svc.Refresh(s.ctx, refreshIn(rt))
	wantErr := auth_service.ErrSessionRevoked
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestRefresh_Expired() {
	now := time.Now()
	rt := "refresh1"

	s.st.EXPECT().
		GetSessionByRefreshHash(mock.Anything, sha256b(rt)).
		Return(&models.Session{
			ID:        10,
			UserID:    1,
			ExpiresAt: now.Add(-time.Second),
			RevokedAt: nil,
		}, nil)
	_, gotErr := s.svc.Refresh(s.ctx, refreshIn(rt))
	wantErr := auth_service.ErrSessionExpired
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestRefresh_RevokeSessionError() {
	now := time.Now()
	rt := "refresh1"
	wantErr := errors.New("revoke fail")

	s.st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b(rt)).
		Return(&models.Session{ID: 10, UserID: 1, ExpiresAt: now.Add(time.Hour)}, nil)
	s.st.EXPECT().RevokeSessionByID(mock.Anything, uint64(10), mock.Anything).Return(wantErr)

	_, gotErr := s.svc.Refresh(s.ctx, refreshIn(rt))
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestRefresh_Success_Rotates() {
	now := time.Now()
	rt := "refresh1"

	s.st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b(rt)).
		Return(&models.Session{ID: 10, UserID: 1, ExpiresAt: now.Add(time.Hour)}, nil)
	s.st.EXPECT().RevokeSessionByID(mock.Anything, uint64(10), mock.Anything).Return(nil)
	s.st.EXPECT().CreateSession(mock.Anything, uint64(1), mock.Anything, mock.Anything, "ua", "ip").
		Return(uint64(11), nil)

	got, gotErr := s.svc.Refresh(s.ctx, refreshIn(rt))
	assert.NilError(s.T(), gotErr)

	wantAccessNonEmpty := true
	wantRefreshNonEmpty := true
	assert.Equal(s.T(), got.AccessToken != "", wantAccessNonEmpty)
	assert.Equal(s.T(), got.RefreshToken != "", wantRefreshNonEmpty)
}

func (s *AuthServiceSuite) TestLogout_InvalidArgs() {
	gotErr := s.svc.Logout(s.ctx, "")
	wantErr := auth_service.ErrInvalidArgument
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestLogout_SessionNotFound() {
	s.st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("tok")).Return(nil, auth_storage.ErrSessionNotFound)

	gotErr := s.svc.Logout(s.ctx, "tok")
	wantErr := auth_service.ErrInvalidRefreshToken
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestLogout_GetSessionOtherError() {
	wantErr := errors.New("db fail")
	s.st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("tok")).Return(nil, wantErr)

	gotErr := s.svc.Logout(s.ctx, "tok")
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestLogout_RevokeError() {
	wantErr := errors.New("revoke fail")
	s.st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("refresh1")).
		Return(&models.Session{ID: 10, UserID: 1, ExpiresAt: time.Now().Add(time.Hour)}, nil)
	s.st.EXPECT().RevokeSessionByID(mock.Anything, uint64(10), mock.Anything).Return(wantErr)

	gotErr := s.svc.Logout(s.ctx, "refresh1")
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestLogout_Success() {
	s.st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("refresh1")).
		Return(&models.Session{ID: 10, UserID: 1, ExpiresAt: time.Now().Add(time.Hour)}, nil)
	s.st.EXPECT().RevokeSessionByID(mock.Anything, uint64(10), mock.Anything).
		Return(nil)

	gotErr := s.svc.Logout(s.ctx, "refresh1")
	assert.NilError(s.T(), gotErr)
}

func (s *AuthServiceSuite) TestLogoutAll_InvalidArgs() {
	gotErr := s.svc.LogoutAll(s.ctx, "")
	wantErr := auth_service.ErrInvalidArgument
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestLogoutAll_SessionNotFound() {
	s.st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("tok")).Return(nil, auth_storage.ErrSessionNotFound)

	gotErr := s.svc.LogoutAll(s.ctx, "tok")
	wantErr := auth_service.ErrInvalidRefreshToken
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestLogoutAll_GetSessionOtherError() {
	wantErr := errors.New("db fail")
	s.st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("tok")).Return(nil, wantErr)

	gotErr := s.svc.LogoutAll(s.ctx, "tok")
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestLogoutAll_RevokeAllError() {
	wantErr := errors.New("revoke all fail")
	s.st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("refresh1")).
		Return(&models.Session{ID: 10, UserID: 1, ExpiresAt: time.Now().Add(time.Hour)}, nil)
	s.st.EXPECT().RevokeAllSessionsByUserID(mock.Anything, uint64(1), mock.Anything).Return(wantErr)

	gotErr := s.svc.LogoutAll(s.ctx, "refresh1")
	assert.ErrorIs(s.T(), gotErr, wantErr)
}

func (s *AuthServiceSuite) TestLogoutAll_Success() {
	s.st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("refresh1")).
		Return(&models.Session{ID: 10, UserID: 1, ExpiresAt: time.Now().Add(time.Hour)}, nil)
	s.st.EXPECT().RevokeAllSessionsByUserID(mock.Anything, uint64(1), mock.Anything).
		Return(nil)

	gotErr := s.svc.LogoutAll(s.ctx, "refresh1")
	assert.NilError(s.T(), gotErr)
}

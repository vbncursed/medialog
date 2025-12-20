package authService_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/vbncursed/medialog/auth-service/internal/models"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService"
	pguserstorage "github.com/vbncursed/medialog/auth-service/internal/storage/pgUserStorage"
	"gotest.tools/v3/assert"
)

func TestAuthService_Refresh_InvalidArgs(t *testing.T) {
	svc, _ := setup(t)
	_, gotErr := svc.Refresh(bg(), models.RefreshInput{RefreshToken: ""})
	wantErr := authService.ErrInvalidArgument
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_Refresh_SessionNotFound(t *testing.T) {
	svc, st := setup(t)
	st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("tok")).Return(nil, pguserstorage.ErrSessionNotFound)
	_, gotErr := svc.Refresh(bg(), refreshIn("tok"))
	wantErr := authService.ErrInvalidRefreshToken
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_Refresh_GetSessionOtherError(t *testing.T) {
	wantErr := errors.New("db fail")
	svc, st := setup(t)
	st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("tok")).Return(nil, wantErr)

	_, gotErr := svc.Refresh(bg(), refreshIn("tok"))
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_Refresh_Revoked(t *testing.T) {
	now := time.Now()
	rt := "refresh1"
	revokedAt := now.Add(-time.Minute)

	svc, st := setup(t)
	st.EXPECT().
		GetSessionByRefreshHash(mock.Anything, sha256b(rt)).
		Return(&models.Session{
			ID:        10,
			UserID:    1,
			ExpiresAt: now.Add(time.Hour),
			RevokedAt: &revokedAt,
		}, nil)
	_, gotErr := svc.Refresh(bg(), refreshIn(rt))
	wantErr := authService.ErrSessionRevoked
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_Refresh_Expired(t *testing.T) {
	now := time.Now()
	rt := "refresh1"

	svc, st := setup(t)
	st.EXPECT().
		GetSessionByRefreshHash(mock.Anything, sha256b(rt)).
		Return(&models.Session{
			ID:        10,
			UserID:    1,
			ExpiresAt: now.Add(-time.Second),
			RevokedAt: nil,
		}, nil)
	_, gotErr := svc.Refresh(bg(), refreshIn(rt))
	wantErr := authService.ErrSessionExpired
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_Refresh_RevokeSessionError(t *testing.T) {
	now := time.Now()
	rt := "refresh1"
	wantErr := errors.New("revoke fail")

	svc, st := setup(t)
	st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b(rt)).
		Return(&models.Session{ID: 10, UserID: 1, ExpiresAt: now.Add(time.Hour)}, nil)
	st.EXPECT().RevokeSessionByID(mock.Anything, uint64(10), mock.Anything).Return(wantErr)

	_, gotErr := svc.Refresh(bg(), refreshIn(rt))
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_Refresh_Success_Rotates(t *testing.T) {
	now := time.Now()
	rt := "refresh1"

	svc, st := setup(t)
	st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b(rt)).
		Return(&models.Session{ID: 10, UserID: 1, ExpiresAt: now.Add(time.Hour)}, nil)
	st.EXPECT().RevokeSessionByID(mock.Anything, uint64(10), mock.Anything).Return(nil)
	st.EXPECT().CreateSession(mock.Anything, uint64(1), mock.Anything, mock.Anything, "ua", "ip").
		Return(uint64(11), nil)

	got, gotErr := svc.Refresh(bg(), refreshIn(rt))
	assert.NilError(t, gotErr)

	wantAccessNonEmpty := true
	wantRefreshNonEmpty := true
	assert.Equal(t, got.AccessToken != "", wantAccessNonEmpty)
	assert.Equal(t, got.RefreshToken != "", wantRefreshNonEmpty)
}

func TestAuthService_Logout_InvalidArgs(t *testing.T) {
	svc, _ := setup(t)
	gotErr := svc.Logout(bg(), "")
	wantErr := authService.ErrInvalidArgument
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_Logout_SessionNotFound(t *testing.T) {
	svc, st := setup(t)
	st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("tok")).Return(nil, pguserstorage.ErrSessionNotFound)

	gotErr := svc.Logout(bg(), "tok")
	wantErr := authService.ErrInvalidRefreshToken
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_Logout_GetSessionOtherError(t *testing.T) {
	wantErr := errors.New("db fail")
	svc, st := setup(t)
	st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("tok")).Return(nil, wantErr)

	gotErr := svc.Logout(bg(), "tok")
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_Logout_RevokeError(t *testing.T) {
	wantErr := errors.New("revoke fail")
	svc, st := setup(t)
	st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("refresh1")).
		Return(&models.Session{ID: 10, UserID: 1, ExpiresAt: time.Now().Add(time.Hour)}, nil)
	st.EXPECT().RevokeSessionByID(mock.Anything, uint64(10), mock.Anything).Return(wantErr)

	gotErr := svc.Logout(bg(), "refresh1")
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_Logout_Success(t *testing.T) {
	svc, st := setup(t)
	st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("refresh1")).
		Return(&models.Session{ID: 10, UserID: 1, ExpiresAt: time.Now().Add(time.Hour)}, nil)
	st.EXPECT().RevokeSessionByID(mock.Anything, uint64(10), mock.Anything).
		Return(nil)

	gotErr := svc.Logout(bg(), "refresh1")
	assert.NilError(t, gotErr)
}

func TestAuthService_LogoutAll_InvalidArgs(t *testing.T) {
	svc, _ := setup(t)
	gotErr := svc.LogoutAll(bg(), "")
	wantErr := authService.ErrInvalidArgument
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_LogoutAll_SessionNotFound(t *testing.T) {
	svc, st := setup(t)
	st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("tok")).Return(nil, pguserstorage.ErrSessionNotFound)

	gotErr := svc.LogoutAll(bg(), "tok")
	wantErr := authService.ErrInvalidRefreshToken
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_LogoutAll_GetSessionOtherError(t *testing.T) {
	wantErr := errors.New("db fail")
	svc, st := setup(t)
	st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("tok")).Return(nil, wantErr)

	gotErr := svc.LogoutAll(bg(), "tok")
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_LogoutAll_RevokeAllError(t *testing.T) {
	wantErr := errors.New("revoke all fail")
	svc, st := setup(t)
	st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("refresh1")).
		Return(&models.Session{ID: 10, UserID: 1, ExpiresAt: time.Now().Add(time.Hour)}, nil)
	st.EXPECT().RevokeAllSessionsByUserID(mock.Anything, uint64(1), mock.Anything).Return(wantErr)

	gotErr := svc.LogoutAll(bg(), "refresh1")
	assert.ErrorIs(t, gotErr, wantErr)
}

func TestAuthService_LogoutAll_Success(t *testing.T) {
	svc, st := setup(t)
	st.EXPECT().GetSessionByRefreshHash(mock.Anything, sha256b("refresh1")).
		Return(&models.Session{ID: 10, UserID: 1, ExpiresAt: time.Now().Add(time.Hour)}, nil)
	st.EXPECT().RevokeAllSessionsByUserID(mock.Anything, uint64(1), mock.Anything).
		Return(nil)

	gotErr := svc.LogoutAll(bg(), "refresh1")
	assert.NilError(t, gotErr)
}

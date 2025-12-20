package authService

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAuthService_Refresh_InvalidArgs(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	_, err := svc.Refresh(context.Background(), "", "", "127.0.0.1")
	require.ErrorIs(t, err, ErrInvalidArgument)
}

func TestAuthService_Refresh_SessionNotFound(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	_, err := svc.Refresh(context.Background(), "tok", "", "127.0.0.1")
	require.ErrorIs(t, err, ErrInvalidRefreshToken)
}

func TestAuthService_Refresh_TokenToHashError(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	old := tokenToHashFn
	tokenToHashFn = func(_ string) (string, []byte, time.Time, error) {
		return "", nil, time.Time{}, errors.New("hash fail")
	}
	t.Cleanup(func() { tokenToHashFn = old })

	_, err := svc.Refresh(context.Background(), "refresh1", "ua", "ip")
	require.ErrorIs(t, err, ErrInvalidRefreshToken)
}

func TestAuthService_Refresh_GetSessionOtherError(t *testing.T) {
	st := newFakeStorage()
	st.errGetSession = errors.New("db fail")
	svc := NewAuthService(st, "secret", 60, 3600)

	_, err := svc.Refresh(context.Background(), "refresh1", "ua", "ip")
	require.Error(t, err)
}

func TestAuthService_Refresh_Revoked(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	rt := "refresh1"
	h := sha256b(rt)
	now := time.Now()
	_, _ = st.CreateSession(context.Background(), 1, h, now.Add(time.Hour), "ua", "ip")
	sess, _ := st.GetSessionByRefreshHash(context.Background(), h)
	sess.RevokedAt = &now

	_, err := svc.Refresh(context.Background(), rt, "ua", "ip")
	require.ErrorIs(t, err, ErrSessionRevoked)
}

func TestAuthService_Refresh_Expired(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	rt := "refresh1"
	h := sha256b(rt)
	_, _ = st.CreateSession(context.Background(), 1, h, time.Now().Add(-time.Minute), "ua", "ip")

	_, err := svc.Refresh(context.Background(), rt, "ua", "ip")
	require.ErrorIs(t, err, ErrSessionExpired)
}

func TestAuthService_Refresh_Success_Rotates(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	rt := "refresh1"
	h := sha256b(rt)
	_, _ = st.CreateSession(context.Background(), 1, h, time.Now().Add(time.Hour), "ua", "ip")

	res, err := svc.Refresh(context.Background(), rt, "ua", "ip")
	require.NoError(t, err)
	require.NotEmpty(t, res.AccessToken)
	require.NotEmpty(t, res.RefreshToken)

	// старую сессию должны были отозвать
	sess, _ := st.GetSessionByRefreshHash(context.Background(), h)
	require.NotNil(t, sess.RevokedAt)
}

func TestAuthService_Logout_InvalidArgs(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	require.ErrorIs(t, svc.Logout(context.Background(), ""), ErrInvalidArgument)
}

func TestAuthService_Logout_SessionNotFound(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	require.ErrorIs(t, svc.Logout(context.Background(), "tok"), ErrInvalidRefreshToken)
}

func TestAuthService_Logout_TokenToHashError(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	old := tokenToHashFn
	tokenToHashFn = func(_ string) (string, []byte, time.Time, error) {
		return "", nil, time.Time{}, errors.New("hash fail")
	}
	t.Cleanup(func() { tokenToHashFn = old })

	require.ErrorIs(t, svc.Logout(context.Background(), "refresh1"), ErrInvalidRefreshToken)
}

func TestAuthService_Logout_Success(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	rt := "refresh1"
	h := sha256b(rt)
	_, _ = st.CreateSession(context.Background(), 1, h, time.Now().Add(time.Hour), "ua", "ip")

	require.NoError(t, svc.Logout(context.Background(), rt))
	sess, _ := st.GetSessionByRefreshHash(context.Background(), h)
	require.NotNil(t, sess.RevokedAt)
}

func TestAuthService_LogoutAll_InvalidArgs(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	require.ErrorIs(t, svc.LogoutAll(context.Background(), ""), ErrInvalidArgument)
}

func TestAuthService_LogoutAll_SessionNotFound(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	require.ErrorIs(t, svc.LogoutAll(context.Background(), "tok"), ErrInvalidRefreshToken)
}

func TestAuthService_LogoutAll_GetSessionOtherError(t *testing.T) {
	st := newFakeStorage()
	st.errGetSession = errors.New("db fail")
	svc := NewAuthService(st, "secret", 60, 3600)

	require.Error(t, svc.LogoutAll(context.Background(), "refresh1"))
}

func TestAuthService_LogoutAll_TokenToHashError(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	old := tokenToHashFn
	tokenToHashFn = func(_ string) (string, []byte, time.Time, error) {
		return "", nil, time.Time{}, errors.New("hash fail")
	}
	t.Cleanup(func() { tokenToHashFn = old })

	require.ErrorIs(t, svc.LogoutAll(context.Background(), "refresh1"), ErrInvalidRefreshToken)
}

func TestAuthService_LogoutAll_Success(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	rt1 := "refresh1"
	rt2 := "refresh2"
	h1 := sha256b(rt1)
	h2 := sha256b(rt2)
	_, _ = st.CreateSession(context.Background(), 1, h1, time.Now().Add(time.Hour), "ua", "ip")
	_, _ = st.CreateSession(context.Background(), 1, h2, time.Now().Add(time.Hour), "ua", "ip")

	require.NoError(t, svc.LogoutAll(context.Background(), rt1))

	s1, _ := st.GetSessionByRefreshHash(context.Background(), h1)
	s2, _ := st.GetSessionByRefreshHash(context.Background(), h2)
	require.NotNil(t, s1.RevokedAt)
	require.NotNil(t, s2.RevokedAt)
}

func TestAuthService_IssueTokens_CreateSessionError(t *testing.T) {
	st := newFakeStorage()
	st.errCreateSession = errors.New("create session fail")
	svc := NewAuthService(st, "secret", 60, 3600)

	_, err := svc.issueTokens(context.Background(), 1, "ua", "ip")
	require.Error(t, err)
}

func TestAuthService_IssueTokens_AccessTokenError(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	old := newAccessTokenFn
	newAccessTokenFn = func(_ string, _ uint64, _ time.Duration) (string, error) {
		return "", errors.New("sign fail")
	}
	t.Cleanup(func() { newAccessTokenFn = old })

	_, err := svc.issueTokens(context.Background(), 1, "ua", "ip")
	require.Error(t, err)
}

func TestAuthService_IssueTokens_RefreshTokenError(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	old := newRefreshTokenFn
	newRefreshTokenFn = func(_ time.Duration) (string, []byte, time.Time, error) {
		return "", nil, time.Time{}, errors.New("refresh fail")
	}
	t.Cleanup(func() { newRefreshTokenFn = old })

	_, err := svc.issueTokens(context.Background(), 1, "ua", "ip")
	require.Error(t, err)
}

func TestAuthService_IssueTokens_RandError(t *testing.T) {
	st := newFakeStorage()
	svc := NewAuthService(st, "secret", 60, 3600)

	old := randRead
	randRead = func(_ []byte) (int, error) { return 0, errors.New("rand fail") }
	t.Cleanup(func() { randRead = old })

	_, err := svc.issueTokens(context.Background(), 1, "ua", "ip")
	require.Error(t, err)
}

func TestAuthService_Refresh_RevokeSessionError(t *testing.T) {
	st := newFakeStorage()
	st.errRevokeSession = errors.New("revoke fail")
	svc := NewAuthService(st, "secret", 60, 3600)

	rt := "refresh1"
	h := sha256b(rt)
	_, _ = st.CreateSession(context.Background(), 1, h, time.Now().Add(time.Hour), "ua", "ip")

	_, err := svc.Refresh(context.Background(), rt, "ua", "ip")
	require.Error(t, err)
}

func TestAuthService_Logout_GetSessionOtherError(t *testing.T) {
	st := newFakeStorage()
	st.errGetSession = errors.New("db fail")
	svc := NewAuthService(st, "secret", 60, 3600)

	require.Error(t, svc.Logout(context.Background(), "refresh1"))
}

func TestAuthService_Logout_RevokeError(t *testing.T) {
	st := newFakeStorage()
	st.errRevokeSession = errors.New("revoke fail")
	svc := NewAuthService(st, "secret", 60, 3600)

	rt := "refresh1"
	h := sha256b(rt)
	_, _ = st.CreateSession(context.Background(), 1, h, time.Now().Add(time.Hour), "ua", "ip")

	require.Error(t, svc.Logout(context.Background(), rt))
}

func TestAuthService_LogoutAll_RevokeAllError(t *testing.T) {
	st := newFakeStorage()
	st.errRevokeAll = errors.New("revoke all fail")
	svc := NewAuthService(st, "secret", 60, 3600)

	rt := "refresh1"
	h := sha256b(rt)
	_, _ = st.CreateSession(context.Background(), 1, h, time.Now().Add(time.Hour), "ua", "ip")

	require.Error(t, svc.LogoutAll(context.Background(), rt))
}

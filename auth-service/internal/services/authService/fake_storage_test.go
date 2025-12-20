package authService

import (
	"context"
	"crypto/sha256"
	"errors"
	"time"

	"github.com/vbncursed/medialog/auth-service/internal/models"
	"github.com/vbncursed/medialog/auth-service/internal/storage/pgstorage"
)

type fakeStorage struct {
	nextUserID    uint64
	nextSessionID uint64

	usersByEmail map[string]*models.User
	// key: string(refresh_hash)
	sessionsByHash map[string]*models.Session

	// error injection
	errGetUser       error
	errCreateUser    error
	errCreateSession error
	errGetSession    error
	errRevokeSession error
	errRevokeAll     error
}

func newFakeStorage() *fakeStorage {
	return &fakeStorage{
		nextUserID:     1,
		nextSessionID:  1,
		usersByEmail:   make(map[string]*models.User),
		sessionsByHash: make(map[string]*models.Session),
	}
}

func (s *fakeStorage) CreateUser(ctx context.Context, email string, passwordHash string) (uint64, error) {
	if s.errCreateUser != nil {
		return 0, s.errCreateUser
	}
	if _, ok := s.usersByEmail[email]; ok {
		return 0, errors.New("duplicate")
	}
	id := s.nextUserID
	s.nextUserID++
	s.usersByEmail[email] = &models.User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}
	return id, nil
}

func (s *fakeStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if s.errGetUser != nil {
		return nil, s.errGetUser
	}
	u, ok := s.usersByEmail[email]
	if !ok {
		return nil, pgstorage.ErrUserNotFound
	}
	return u, nil
}

func (s *fakeStorage) CreateSession(ctx context.Context, userID uint64, refreshHash []byte, expiresAt time.Time, userAgent, ip string) (uint64, error) {
	if s.errCreateSession != nil {
		return 0, s.errCreateSession
	}
	id := s.nextSessionID
	s.nextSessionID++
	h := string(refreshHash)
	s.sessionsByHash[h] = &models.Session{
		ID:          id,
		UserID:      userID,
		RefreshHash: refreshHash,
		ExpiresAt:   expiresAt,
		CreatedAt:   time.Now(),
		UserAgent:   userAgent,
		IP:          ip,
	}
	return id, nil
}

func (s *fakeStorage) GetSessionByRefreshHash(ctx context.Context, refreshHash []byte) (*models.Session, error) {
	if s.errGetSession != nil {
		return nil, s.errGetSession
	}
	sess, ok := s.sessionsByHash[string(refreshHash)]
	if !ok {
		return nil, pgstorage.ErrSessionNotFound
	}
	return sess, nil
}

func (s *fakeStorage) RevokeSessionByID(ctx context.Context, sessionID uint64, revokedAt time.Time) error {
	if s.errRevokeSession != nil {
		return s.errRevokeSession
	}
	for _, sess := range s.sessionsByHash {
		if sess.ID == sessionID {
			if sess.RevokedAt != nil {
				return pgstorage.ErrSessionNotFound
			}
			sess.RevokedAt = &revokedAt
			return nil
		}
	}
	return pgstorage.ErrSessionNotFound
}

func (s *fakeStorage) RevokeAllSessionsByUserID(ctx context.Context, userID uint64, revokedAt time.Time) error {
	if s.errRevokeAll != nil {
		return s.errRevokeAll
	}
	for _, sess := range s.sessionsByHash {
		if sess.UserID == userID && sess.RevokedAt == nil {
			sess.RevokedAt = &revokedAt
		}
	}
	return nil
}

func sha256b(v string) []byte {
	sum := sha256.Sum256([]byte(v))
	return sum[:]
}

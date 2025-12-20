package authService

import (
	"context"
	"crypto/sha256"
	"errors"
	"strings"
	"time"

	"github.com/vbncursed/medialog/auth-service/internal/models"
	pgstorage "github.com/vbncursed/medialog/auth-service/internal/storage/pgUserStorage"
)

type AuthService struct {
	storage Storage

	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// test hooks for hard-to-trigger branches
var (
	tokenToHashFn     = tokenToHash
	newAccessTokenFn  = newAccessToken
	newRefreshTokenFn = newRefreshToken
)

func NewAuthService(storage Storage, jwtSecret string, accessTTLSeconds, refreshTTLSeconds int64) *AuthService {
	return &AuthService{
		storage:    storage,
		jwtSecret:  jwtSecret,
		accessTTL:  time.Duration(accessTTLSeconds) * time.Second,
		refreshTTL: time.Duration(refreshTTLSeconds) * time.Second,
	}
}

type AuthResult struct {
	UserID       uint64
	AccessToken  string
	RefreshToken string
}

func (s *AuthService) Register(ctx context.Context, in models.RegisterInput) (*AuthResult, error) {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	if !validateEmail(in.Email) || !validatePassword(in.Password) {
		return nil, ErrInvalidArgument
	}

	// Проверяем существование пользователя.
	if _, err := s.storage.GetUserByEmail(ctx, in.Email); err == nil {
		return nil, ErrEmailAlreadyExists
	} else if !errors.Is(err, pgstorage.ErrUserNotFound) {
		return nil, err
	}

	passHash, err := passwordHash(in.Password)
	if err != nil {
		return nil, err
	}

	userID, err := s.storage.CreateUser(ctx, in.Email, passHash)
	if err != nil {
		// Уникальность email может “стрельнуть” гонкой.
		return nil, ErrEmailAlreadyExists
	}

	return s.issueTokens(ctx, userID, in.UserAgent, in.IP)
}

func (s *AuthService) Login(ctx context.Context, in models.LoginInput) (*AuthResult, error) {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	if !validateEmail(in.Email) || !validatePassword(in.Password) {
		return nil, ErrInvalidArgument
	}

	u, err := s.storage.GetUserByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, pgstorage.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !comparePassword(u.PasswordHash, in.Password) {
		return nil, ErrInvalidCredentials
	}

	return s.issueTokens(ctx, u.ID, in.UserAgent, in.IP)
}

func (s *AuthService) Refresh(ctx context.Context, in models.RefreshInput) (*AuthResult, error) {
	in.RefreshToken = strings.TrimSpace(in.RefreshToken)
	if in.RefreshToken == "" {
		return nil, ErrInvalidArgument
	}

	// Хешируем refresh token (в БД хранится только hash).
	_, refreshHash, _, err := tokenToHashFn(in.RefreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	sess, err := s.storage.GetSessionByRefreshHash(ctx, refreshHash)
	if err != nil {
		if errors.Is(err, pgstorage.ErrSessionNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, err
	}
	if sess.RevokedAt != nil {
		return nil, ErrSessionRevoked
	}
	if time.Now().After(sess.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	// Ротация: ревокаем старую сессию, создаём новую.
	now := time.Now()
	if err := s.storage.RevokeSessionByID(ctx, sess.ID, now); err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, sess.UserID, in.UserAgent, in.IP)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return ErrInvalidArgument
	}

	_, refreshHash, _, err := tokenToHashFn(refreshToken)
	if err != nil {
		return ErrInvalidRefreshToken
	}

	sess, err := s.storage.GetSessionByRefreshHash(ctx, refreshHash)
	if err != nil {
		if errors.Is(err, pgstorage.ErrSessionNotFound) {
			return ErrInvalidRefreshToken
		}
		return err
	}

	now := time.Now()
	if err := s.storage.RevokeSessionByID(ctx, sess.ID, now); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) LogoutAll(ctx context.Context, refreshToken string) error {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return ErrInvalidArgument
	}

	_, refreshHash, _, err := tokenToHashFn(refreshToken)
	if err != nil {
		return ErrInvalidRefreshToken
	}

	sess, err := s.storage.GetSessionByRefreshHash(ctx, refreshHash)
	if err != nil {
		if errors.Is(err, pgstorage.ErrSessionNotFound) {
			return ErrInvalidRefreshToken
		}
		return err
	}

	now := time.Now()
	if err := s.storage.RevokeAllSessionsByUserID(ctx, sess.UserID, now); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) issueTokens(ctx context.Context, userID uint64, userAgent, ip string) (*AuthResult, error) {
	access, err := newAccessTokenFn(s.jwtSecret, userID, s.accessTTL)
	if err != nil {
		return nil, err
	}

	refresh, refreshHash, refreshExp, err := newRefreshTokenFn(s.refreshTTL)
	if err != nil {
		return nil, err
	}

	_, err = s.storage.CreateSession(ctx, userID, refreshHash, refreshExp, userAgent, ip)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		UserID:       userID,
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func tokenToHash(refreshToken string) (token string, hash []byte, exp time.Time, err error) {
	// refresh токен не содержит exp внутри (opaque).
	// Это helper только для хеширования единообразно.
	sum := sha256.Sum256([]byte(refreshToken))
	return refreshToken, sum[:], time.Time{}, nil
}

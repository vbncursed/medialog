package auth_service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/vbncursed/medialog/auth-service/internal/models"
	"github.com/vbncursed/medialog/auth-service/internal/storage/auth_storage"
)

func (s *AuthService) Refresh(ctx context.Context, in models.RefreshInput) (*AuthInfo, error) {
	in.RefreshToken = strings.TrimSpace(in.RefreshToken)
	if in.RefreshToken == "" {
		return nil, ErrInvalidArgument
	}

	_, refreshHash, _, err := tokenToHashFn(in.RefreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	sess, err := s.authStorage.GetSessionByRefreshHash(ctx, refreshHash)
	if err != nil {
		if errors.Is(err, auth_storage.ErrSessionNotFound) {
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

	now := time.Now()
	if err := s.authStorage.RevokeSessionByID(ctx, sess.ID, now); err != nil {
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

	sess, err := s.authStorage.GetSessionByRefreshHash(ctx, refreshHash)
	if err != nil {
		if errors.Is(err, auth_storage.ErrSessionNotFound) {
			return ErrInvalidRefreshToken
		}
		return err
	}

	now := time.Now()
	if err := s.authStorage.RevokeSessionByID(ctx, sess.ID, now); err != nil {
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

	sess, err := s.authStorage.GetSessionByRefreshHash(ctx, refreshHash)
	if err != nil {
		if errors.Is(err, auth_storage.ErrSessionNotFound) {
			return ErrInvalidRefreshToken
		}
		return err
	}

	now := time.Now()
	if err := s.authStorage.RevokeAllSessionsByUserID(ctx, sess.UserID, now); err != nil {
		return err
	}

	return nil
}

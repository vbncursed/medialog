package auth_service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/vbncursed/medialog/auth_service/internal/models"
)

var (
	ErrSessionNotFound = errors.New("session not found")
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

	sess, err := s.sessionStorage.GetSessionByRefreshHash(ctx, refreshHash)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
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

	if err := s.sessionStorage.RevokeSessionByRefreshHash(ctx, refreshHash); err != nil {
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

	_, err = s.sessionStorage.GetSessionByRefreshHash(ctx, refreshHash)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return ErrInvalidRefreshToken
		}
		return err
	}

	if err := s.sessionStorage.RevokeSessionByRefreshHash(ctx, refreshHash); err != nil {
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

	sess, err := s.sessionStorage.GetSessionByRefreshHash(ctx, refreshHash)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return ErrInvalidRefreshToken
		}
		return err
	}

	if err := s.sessionStorage.RevokeAllSessionsByUserID(ctx, sess.UserID); err != nil {
		return err
	}

	return nil
}

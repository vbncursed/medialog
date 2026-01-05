package auth_service

import (
	"context"
	"time"

	"github.com/vbncursed/medialog/auth/internal/models"
)

func (s *AuthService) Refresh(ctx context.Context, in models.RefreshInput) (*models.AuthInfo, error) {
	if in.RefreshToken == "" {
		return nil, ErrInvalidArgument
	}

	refreshHash := tokenToHash(in.RefreshToken)

	sess, err := s.sessionStorage.GetSessionByRefreshHash(ctx, refreshHash)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	if time.Now().After(sess.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	user, err := s.authStorage.GetUserByID(ctx, sess.UserID)
	if err != nil {
		return nil, err
	}

	err = s.sessionStorage.RevokeSessionByRefreshHash(ctx, refreshHash)
	if err != nil {
		return nil, err
	}

	return s.issueTokens(ctx, user.ID, user.Role, in.UserAgent, in.IP)
}

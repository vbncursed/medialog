package authService

import "context"

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

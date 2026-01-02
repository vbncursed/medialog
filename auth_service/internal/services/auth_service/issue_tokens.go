package auth_service

import "context"

func (s *AuthService) issueTokens(ctx context.Context, userID uint64, userAgent, ip string) (*AuthInfo, error) {
	user, err := s.authStorage.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	access, err := newAccessTokenFn(s.jwtSecret, userID, user.Role, s.accessTTL)
	if err != nil {
		return nil, err
	}

	refresh, refreshHash, refreshExp, err := newRefreshTokenFn(s.refreshTTL)
	if err != nil {
		return nil, err
	}

	err = s.sessionStorage.CreateSession(ctx, userID, refreshHash, refreshExp, userAgent, ip)
	if err != nil {
		return nil, err
	}

	return &AuthInfo{
		UserID:       userID,
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

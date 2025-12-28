package auth_service_api

import (
	"context"
	"errors"

	domain "github.com/vbncursed/medialog/auth-service/internal/models"
	"github.com/vbncursed/medialog/auth-service/internal/pb/models"
	"github.com/vbncursed/medialog/auth-service/internal/services/auth_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *AuthServiceAPI) Refresh(ctx context.Context, req *models.RefreshRequest) (*models.AuthResponse, error) {
	ua, ip := clientMeta(ctx)

	if !a.refreshLimiter.Allow(ctx, ip) {
		return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
	}

	res, err := a.authService.Refresh(ctx, domain.RefreshInput{
		RefreshToken: req.GetRefreshToken(),
		UserAgent:    ua,
		IP:           ip,
	})
	if err != nil {
		switch {
		case errors.Is(err, auth_service.ErrInvalidArgument):
			return nil, status.Error(codes.InvalidArgument, "refresh_token required")
		case errors.Is(err, auth_service.ErrInvalidRefreshToken),
			errors.Is(err, auth_service.ErrSessionExpired),
			errors.Is(err, auth_service.ErrSessionRevoked):
			return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &models.AuthResponse{
		UserId:       res.UserID,
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

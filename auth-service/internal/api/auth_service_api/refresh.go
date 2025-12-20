package auth_service_api

import (
	"context"
	"errors"

	"github.com/vbncursed/medialog/auth-service/internal/pb/models"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *AuthServiceAPI) Refresh(ctx context.Context, req *models.RefreshRequest) (*models.AuthResponse, error) {
	ua, ip := clientMeta(ctx)

	res, err := a.authService.Refresh(ctx, req.GetRefreshToken(), ua, ip)
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrInvalidArgument):
			return nil, status.Error(codes.InvalidArgument, "refresh_token required")
		case errors.Is(err, authService.ErrInvalidRefreshToken),
			errors.Is(err, authService.ErrSessionExpired),
			errors.Is(err, authService.ErrSessionRevoked):
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

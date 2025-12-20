package auth_service_api

import (
	"context"
	"errors"

	"github.com/vbncursed/medialog/auth-service/internal/pb/models"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *AuthServiceAPI) Logout(ctx context.Context, req *models.LogoutRequest) (*models.LogoutResponse, error) {
	err := a.authService.Logout(ctx, req.GetRefreshToken())
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrInvalidArgument):
			return nil, status.Error(codes.InvalidArgument, "refresh_token required")
		case errors.Is(err, authService.ErrInvalidRefreshToken):
			return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &models.LogoutResponse{}, nil
}

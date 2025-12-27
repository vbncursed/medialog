package auth_service_api

import (
	"context"
	"errors"

	"github.com/vbncursed/medialog/auth-service/internal/pb/models"
	"github.com/vbncursed/medialog/auth-service/internal/services/auth_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *AuthServiceAPI) LogoutAll(ctx context.Context, req *models.LogoutAllRequest) (*models.LogoutResponse, error) {
	err := a.authService.LogoutAll(ctx, req.GetRefreshToken())
	if err != nil {
		switch {
		case errors.Is(err, auth_service.ErrInvalidArgument):
			return nil, status.Error(codes.InvalidArgument, "refresh_token required")
		case errors.Is(err, auth_service.ErrInvalidRefreshToken):
			return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &models.LogoutResponse{}, nil
}

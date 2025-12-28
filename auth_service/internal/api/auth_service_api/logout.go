package auth_service_api

import (
	"context"
	"errors"

	"github.com/vbncursed/medialog/auth_service/internal/pb/models"
	"github.com/vbncursed/medialog/auth_service/internal/services/auth_service"
	"google.golang.org/grpc/codes"
)

func (a *AuthServiceAPI) Logout(ctx context.Context, req *models.LogoutRequest) (*models.LogoutResponse, error) {
	err := a.authService.Logout(ctx, req.GetRefreshToken())
	if err != nil {
		switch {
		case errors.Is(err, auth_service.ErrInvalidArgument):
			return nil, newFieldError(codes.InvalidArgument, ErrCodeMissingField, "refresh_token", "Refresh token is required.")
		case errors.Is(err, auth_service.ErrInvalidRefreshToken):
			return nil, newError(codes.Unauthenticated, ErrCodeInvalidToken, "Invalid refresh token.")
		default:
			return nil, newError(codes.Internal, ErrCodeInternal, "An internal error occurred. Please try again later.")
		}
	}

	return &models.LogoutResponse{}, nil
}

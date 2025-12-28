package auth_service_api

import (
	"context"
	"errors"

	domain "github.com/vbncursed/medialog/auth_service/internal/models"
	"github.com/vbncursed/medialog/auth_service/internal/pb/models"
	"github.com/vbncursed/medialog/auth_service/internal/services/auth_service"
	"google.golang.org/grpc/codes"
)

func (a *AuthServiceAPI) Refresh(ctx context.Context, req *models.RefreshRequest) (*models.AuthResponse, error) {
	ua, ip := clientMeta(ctx)

	if !a.refreshLimiter.Allow(ctx, ip) {
		return nil, newError(codes.ResourceExhausted, ErrCodeRateLimitExceeded, "Too many refresh attempts. Please try again later.")
	}

	res, err := a.authService.Refresh(ctx, domain.RefreshInput{
		RefreshToken: req.GetRefreshToken(),
		UserAgent:    ua,
		IP:           ip,
	})
	if err != nil {
		switch {
		case errors.Is(err, auth_service.ErrInvalidArgument):
			return nil, newFieldError(codes.InvalidArgument, ErrCodeMissingField, "refresh_token", "Refresh token is required.")
		case errors.Is(err, auth_service.ErrInvalidRefreshToken):
			return nil, newError(codes.Unauthenticated, ErrCodeInvalidToken, "Invalid refresh token.")
		case errors.Is(err, auth_service.ErrSessionExpired):
			return nil, newError(codes.Unauthenticated, ErrCodeSessionExpired, "Session has expired. Please log in again.")
		case errors.Is(err, auth_service.ErrSessionRevoked):
			return nil, newError(codes.Unauthenticated, ErrCodeSessionRevoked, "Session has been revoked. Please log in again.")
		default:
			if isDatabaseError(err) {
				return nil, newError(codes.Unavailable, ErrCodeServiceUnavailable, "Service temporarily unavailable. Please try again later.")
			}
			return nil, newError(codes.Internal, ErrCodeInternal, "An internal error occurred. Please try again later.")
		}
	}

	return &models.AuthResponse{
		UserId:       res.UserID,
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

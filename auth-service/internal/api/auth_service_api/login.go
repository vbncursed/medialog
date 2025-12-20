package auth_service_api

import (
	"context"
	"errors"

	"github.com/vbncursed/medialog/auth-service/internal/pb/models"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *AuthServiceAPI) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	ua, ip := clientMeta(ctx)

	if !a.loginLimiter.Allow(normalizeRateLimitKey(ip)) {
		return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
	}

	res, err := a.authService.Login(ctx, req.GetEmail(), req.GetPassword(), ua, ip)
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrInvalidArgument):
			return nil, status.Error(codes.InvalidArgument, "invalid email/password")
		case errors.Is(err, authService.ErrInvalidCredentials):
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
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

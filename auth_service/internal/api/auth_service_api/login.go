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

func (a *AuthServiceAPI) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	ua, ip := clientMeta(ctx)

	if !a.loginLimiter.Allow(ctx, normalizeRateLimitKey(ip)) {
		return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
	}

	res, err := a.authService.Login(ctx, domain.LoginInput{
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		UserAgent: ua,
		IP:        ip,
	})
	if err != nil {
		switch {
		case errors.Is(err, auth_service.ErrInvalidArgument):
			return nil, status.Error(codes.InvalidArgument, "invalid email/password")
		case errors.Is(err, auth_service.ErrInvalidCredentials):
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

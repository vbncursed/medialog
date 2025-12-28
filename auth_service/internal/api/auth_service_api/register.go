package auth_service_api

import (
	"context"
	"errors"

	domain "github.com/vbncursed/medialog/auth_service/internal/models"
	"github.com/vbncursed/medialog/auth_service/internal/pb/models"
	"github.com/vbncursed/medialog/auth_service/internal/services/auth_service"
	"google.golang.org/grpc/codes"
)

func (a *AuthServiceAPI) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	ua, ip := clientMeta(ctx)

	if !a.registerLimiter.Allow(ctx, ip) {
		return nil, newError(codes.ResourceExhausted, ErrCodeRateLimitExceeded, "Too many registration attempts. Please try again later.")
	}

	res, err := a.authService.Register(ctx, domain.RegisterInput{
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		UserAgent: ua,
		IP:        ip,
	})
	if err != nil {
		switch {
		case errors.Is(err, auth_service.ErrInvalidEmail):
			return nil, newFieldError(codes.InvalidArgument, ErrCodeInvalidEmail, "email", "Invalid email format.")
		case errors.Is(err, auth_service.ErrInvalidPassword):
			return nil, newFieldError(codes.InvalidArgument, ErrCodeInvalidPassword, "password", "Password must be at least 8 characters with uppercase, lowercase, and digit.")
		case errors.Is(err, auth_service.ErrInvalidArgument):
			return nil, newError(codes.InvalidArgument, ErrCodeInvalidInput, "Invalid email or password format.")
		case errors.Is(err, auth_service.ErrEmailAlreadyExists):
			return nil, newError(codes.AlreadyExists, ErrCodeEmailAlreadyExists, "An account with this email already exists.")
		default:
			return nil, newError(codes.Internal, ErrCodeInternal, "An internal error occurred. Please try again later.")
		}
	}

	return &models.AuthResponse{
		UserId:       res.UserID,
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

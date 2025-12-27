package auth_service

import (
	"context"
	"errors"

	"github.com/vbncursed/medialog/auth-service/internal/models"
	"github.com/vbncursed/medialog/auth-service/internal/storage/auth_storage"
)

func (s *AuthService) Register(ctx context.Context, in models.RegisterInput) (*AuthInfo, error) {
	var err error
	in, err = normalizeAndValidateAuthInput(in)
	if err != nil {
		return nil, err
	}

	// Проверяем существование пользователя.
	_, err = s.authStorage.GetUserByEmail(ctx, in.Email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, auth_storage.ErrUserNotFound) {
		return nil, err
	}

	passHash, err := passwordHash(in.Password)
	if err != nil {
		return nil, err
	}

	userID, err := s.authStorage.CreateUser(ctx, in.Email, passHash)
	if err != nil {
		return nil, ErrEmailAlreadyExists
	}

	return s.issueTokens(ctx, userID, in.UserAgent, in.IP)
}

func (s *AuthService) Login(ctx context.Context, in models.LoginInput) (*AuthInfo, error) {
	var err error
	in, err = normalizeAndValidateAuthInput(in)
	if err != nil {
		return nil, err
	}

	u, err := s.authStorage.GetUserByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, auth_storage.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !comparePassword(u.PasswordHash, in.Password) {
		return nil, ErrInvalidCredentials
	}

	return s.issueTokens(ctx, u.ID, in.UserAgent, in.IP)
}

package authService

import (
	"context"
	"errors"

	"github.com/vbncursed/medialog/auth-service/internal/models"
	pguserstorage "github.com/vbncursed/medialog/auth-service/internal/storage/pgUserStorage"
)

func (s *AuthService) Register(ctx context.Context, in models.RegisterInput) (*AuthInfo, error) {
	var err error
	in, err = normalizeAndValidateAuthInput(in)
	if err != nil {
		return nil, err
	}

	// Проверяем существование пользователя.
	if _, err := s.storage.GetUserByEmail(ctx, in.Email); err == nil {
		return nil, ErrEmailAlreadyExists
	} else if !errors.Is(err, pguserstorage.ErrUserNotFound) {
		return nil, err
	}

	passHash, err := passwordHash(in.Password)
	if err != nil {
		return nil, err
	}

	userID, err := s.storage.CreateUser(ctx, in.Email, passHash)
	if err != nil {
		// Уникальность email может “стрельнуть” гонкой.
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

	u, err := s.storage.GetUserByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, pguserstorage.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !comparePassword(u.PasswordHash, in.Password) {
		return nil, ErrInvalidCredentials
	}

	return s.issueTokens(ctx, u.ID, in.UserAgent, in.IP)
}

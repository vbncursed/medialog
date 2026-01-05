package auth_service

import (
	"context"

	"github.com/vbncursed/medialog/auth/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Login(ctx context.Context, in models.LoginInput) (*models.AuthInfo, error) {
	if err := validateAuthInput(in.Email, in.Password); err != nil {
		return nil, err
	}

	user, err := s.authStorage.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.issueTokens(ctx, user.ID, user.Role, in.UserAgent, in.IP)
}

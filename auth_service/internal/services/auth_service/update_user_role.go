package auth_service

import (
	"context"
	"errors"

	"github.com/vbncursed/medialog/auth_service/internal/models"
)

var (
	ErrInvalidRole         = errors.New("invalid role")
	ErrPermissionDenied    = errors.New("permission denied")
	ErrCannotChangeOwnRole = errors.New("cannot change own role")
)

func (s *AuthService) UpdateUserRole(ctx context.Context, adminUserID uint64, targetUserID uint64, newRole string) error {
	if newRole != models.RoleUser && newRole != models.RoleAdmin {
		return ErrInvalidRole
	}

	if adminUserID == targetUserID {
		return ErrCannotChangeOwnRole
	}

	admin, err := s.authStorage.GetUserByID(ctx, adminUserID)
	if err != nil {
		return err
	}

	if admin.Role != models.RoleAdmin {
		return ErrPermissionDenied
	}

	return s.authStorage.UpdateUserRole(ctx, targetUserID, newRole)
}

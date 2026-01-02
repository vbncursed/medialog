package auth_service_api

import (
	"context"
	"errors"

	"github.com/vbncursed/medialog/auth_service/internal/models"
	pb_models "github.com/vbncursed/medialog/auth_service/internal/pb/models"
	"github.com/vbncursed/medialog/auth_service/internal/services/auth_service"
	"google.golang.org/grpc/codes"
)

func (a *AuthServiceAPI) UpdateUserRole(ctx context.Context, req *pb_models.UpdateUserRoleRequest) (*pb_models.UpdateUserRoleResponse, error) {
	// Извлекаем user_id администратора из JWT токена
	adminUserID, err := a.getUserIDFromContext(ctx, a.jwtSecret)
	if err != nil {
		return nil, newError(codes.Unauthenticated, ErrCodeUnauthorized, "Authentication required. Invalid or missing JWT token.")
	}

	targetUserID := req.GetUserId()
	newRole := req.GetRole()

	// Валидация роли
	if newRole != models.RoleUser && newRole != models.RoleAdmin {
		return nil, newFieldError(codes.InvalidArgument, ErrCodeInvalidInput, "role", "Role must be 'user' or 'admin'.")
	}

	err = a.authService.UpdateUserRole(ctx, adminUserID, targetUserID, newRole)
	if err != nil {
		switch {
		case errors.Is(err, auth_service.ErrInvalidRole):
			return nil, newFieldError(codes.InvalidArgument, ErrCodeInvalidInput, "role", "Invalid role. Must be 'user' or 'admin'.")
		case errors.Is(err, auth_service.ErrPermissionDenied):
			return nil, newError(codes.PermissionDenied, ErrCodeUnauthorized, "Only administrators can change user roles.")
		case errors.Is(err, auth_service.ErrCannotChangeOwnRole):
			return nil, newError(codes.InvalidArgument, ErrCodeInvalidInput, "Cannot change your own role.")
		default:
			if isDatabaseError(err) {
				return nil, newError(codes.Unavailable, ErrCodeServiceUnavailable, "Service temporarily unavailable. Please try again later.")
			}
			return nil, newError(codes.Internal, ErrCodeInternal, "An internal error occurred. Please try again later.")
		}
	}

	return &pb_models.UpdateUserRoleResponse{}, nil
}


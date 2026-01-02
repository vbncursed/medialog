package auth_service_api

import (
	"context"

	"github.com/vbncursed/medialog/auth_service/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RequireRole(jwtSecret string, requiredRole string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		role, err := extractRoleFromContext(ctx, jwtSecret)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "Authentication required. Invalid or missing JWT token.")
		}

		if !hasRequiredRole(role, requiredRole) {
			return nil, status.Error(codes.PermissionDenied, "Insufficient permissions. Required role: "+requiredRole)
		}

		return handler(ctx, req)
	}
}

func RequireAnyRole(jwtSecret string, requiredRoles ...string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		role, err := extractRoleFromContext(ctx, jwtSecret)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "Authentication required. Invalid or missing JWT token.")
		}

		for _, requiredRole := range requiredRoles {
			if hasRequiredRole(role, requiredRole) {
				return handler(ctx, req)
			}
		}

		return nil, status.Error(codes.PermissionDenied, "Insufficient permissions.")
	}
}

func extractRoleFromContext(ctx context.Context, jwtSecret string) (string, error) {
	api := &AuthServiceAPI{jwtSecret: jwtSecret}
	return api.getUserRoleFromContext(ctx, jwtSecret)
}

func hasRequiredRole(userRole, requiredRole string) bool {
	if userRole == models.RoleAdmin {
		return true
	}

	if userRole == models.RoleUser {
		return requiredRole == models.RoleUser || requiredRole == models.RoleGuest
	}

	if userRole == models.RoleGuest {
		return requiredRole == models.RoleGuest
	}

	return false
}

func RequireAdmin(jwtSecret string) grpc.UnaryServerInterceptor {
	return RequireRole(jwtSecret, models.RoleAdmin)
}

func RequireUserOrAdmin(jwtSecret string) grpc.UnaryServerInterceptor {
	return RequireAnyRole(jwtSecret, models.RoleUser, models.RoleAdmin)
}

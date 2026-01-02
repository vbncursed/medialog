package library_service_api

import (
	"context"

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
	api := &LibraryServiceAPI{jwtSecret: jwtSecret}
	return api.getUserRoleFromContext(ctx)
}

func hasRequiredRole(userRole, requiredRole string) bool {
	if userRole == "admin" {
		return true
	}

	if userRole == "user" {
		return requiredRole == "user" || requiredRole == "guest"
	}

	if userRole == "guest" {
		return requiredRole == "guest"
	}

	return false
}

func RequireUserOrAdmin(jwtSecret string) grpc.UnaryServerInterceptor {
	return RequireAnyRole(jwtSecret, "user", "admin")
}

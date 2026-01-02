package library_service_api

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RequireRole возвращает gRPC interceptor, который проверяет, что пользователь имеет требуемую роль
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

// RequireAnyRole возвращает interceptor, который проверяет, что пользователь имеет одну из требуемых ролей
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

// extractRoleFromContext извлекает роль из JWT токена в контексте
func extractRoleFromContext(ctx context.Context, jwtSecret string) (string, error) {
	api := &LibraryServiceAPI{jwtSecret: jwtSecret}
	return api.getUserRoleFromContext(ctx)
}

// hasRequiredRole проверяет, имеет ли пользователь требуемую роль
func hasRequiredRole(userRole, requiredRole string) bool {
	// Администратор имеет доступ ко всему
	if userRole == "admin" {
		return true
	}

	// Пользователь имеет доступ только к своим ресурсам и публичным страницам
	if userRole == "user" {
		return requiredRole == "user" || requiredRole == "guest"
	}

	// Гость имеет доступ только к публичным страницам
	if userRole == "guest" {
		return requiredRole == "guest"
	}

	return false
}

// RequireUserOrAdmin проверяет, что пользователь является пользователем или администратором
func RequireUserOrAdmin(jwtSecret string) grpc.UnaryServerInterceptor {
	return RequireAnyRole(jwtSecret, "user", "admin")
}


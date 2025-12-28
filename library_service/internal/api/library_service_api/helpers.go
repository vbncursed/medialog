package library_service_api

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

// getUserIDFromContext извлекает user_id из JWT токена в заголовке Authorization
func (l *LibraryServiceAPI) getUserIDFromContext(ctx context.Context) (uint64, error) {
	// Извлекаем токен из метаданных gRPC
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errors.New("metadata not found in context")
	}

	// Ищем токен в заголовке Authorization
	var tokenString string
	if authHeaders := md.Get("authorization"); len(authHeaders) > 0 {
		authHeader := authHeaders[0]
		// Поддерживаем формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			tokenString = parts[1]
		} else {
			// Если нет префикса Bearer, считаем что весь заголовок - это токен
			tokenString = authHeader
		}
	}

	if tokenString == "" {
		return 0, errors.New("authorization token not found in context")
	}

	// Валидируем и парсим JWT токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(l.jwtSecret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	// Извлекаем user_id из claims (поле "sub" как в auth_service)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	sub, ok := claims["sub"]
	if !ok {
		return 0, errors.New("user_id (sub) not found in token claims")
	}

	// Конвертируем user_id в uint64
	var userID uint64
	switch v := sub.(type) {
	case float64:
		userID = uint64(v)
	case int64:
		userID = uint64(v)
	case uint64:
		userID = v
	case string:
		if _, err := fmt.Sscanf(v, "%d", &userID); err != nil {
			return 0, fmt.Errorf("invalid user_id format in token: %v", v)
		}
	default:
		return 0, fmt.Errorf("unexpected user_id type in token: %T", v)
	}

	if userID == 0 {
		return 0, errors.New("user_id cannot be zero")
	}

	return userID, nil
}

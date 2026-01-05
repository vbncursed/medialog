package auth_service_api

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

// getUserIDFromContext извлекает user_id из JWT токена в заголовке Authorization
func (a *AuthServiceAPI) getUserIDFromContext(ctx context.Context, jwtSecret string) (uint64, error) {
	tokenString, err := extractTokenFromContext(ctx)
	if err != nil {
		return 0, err
	}

	claims, err := parseJWTToken(tokenString, jwtSecret)
	if err != nil {
		return 0, err
	}

	sub, ok := claims["sub"]
	if !ok {
		return 0, errors.New("user_id (sub) not found in token claims")
	}

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

func extractTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata not found in context")
	}

	var tokenString string
	if authHeaders := md.Get("authorization"); len(authHeaders) > 0 {
		authHeader := authHeaders[0]
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			tokenString = parts[1]
		} else {
			tokenString = authHeader
		}
	}

	if tokenString == "" {
		return "", errors.New("authorization token not found in context")
	}

	return tokenString, nil
}

func parseJWTToken(tokenString, jwtSecret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

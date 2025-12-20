package auth_service_api

import (
	"context"
	"net"
	"time"

	"github.com/vbncursed/medialog/auth-service/internal/pb/auth_api"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// AuthServiceAPI реализует grpc AuthServiceServer.
type AuthServiceAPI struct {
	auth_api.UnimplementedAuthServiceServer
	authService     authService.Service
	loginLimiter    *fixedWindowLimiter
	registerLimiter *fixedWindowLimiter
}

func NewAuthServiceAPI(authService authService.Service, loginLimitPerMinute, registerLimitPerMinute int) *AuthServiceAPI {
	return &AuthServiceAPI{
		authService:     authService,
		loginLimiter:    newFixedWindowLimiter(loginLimitPerMinute, time.Minute),
		registerLimiter: newFixedWindowLimiter(registerLimitPerMinute, time.Minute),
	}
}

func clientMeta(ctx context.Context) (userAgent, ip string) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ua := md.Get("user-agent"); len(ua) > 0 {
			userAgent = ua[0]
		}
	}
	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		// peer.Addr может быть *net.TCPAddr, но не гарантировано.
		host, _, err := net.SplitHostPort(p.Addr.String())
		if err == nil {
			ip = host
		} else {
			ip = p.Addr.String()
		}
	}
	return userAgent, ip
}

func normalizeRateLimitKey(ip string) string {
	if ip == "" {
		return "unknown"
	}
	return ip
}

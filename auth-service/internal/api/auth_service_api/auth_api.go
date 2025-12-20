package auth_service_api

import (
	"context"
	"net"
	"net/netip"

	"github.com/vbncursed/medialog/auth-service/internal/pb/auth_api"
	"github.com/vbncursed/medialog/auth-service/internal/services/authService"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// AuthServiceAPI реализует grpc AuthServiceServer.
type AuthServiceAPI struct {
	auth_api.UnimplementedAuthServiceServer
	authService     authService.Service
	loginLimiter    RateLimiter
	registerLimiter RateLimiter
}

type denyAllLimiter struct{}

func (denyAllLimiter) Allow(ctx context.Context, key string) bool { return false }

func NewAuthServiceAPI(authService authService.Service, loginLimiter, registerLimiter RateLimiter) *AuthServiceAPI {
	if loginLimiter == nil {
		loginLimiter = denyAllLimiter{}
	}
	if registerLimiter == nil {
		registerLimiter = denyAllLimiter{}
	}
	return &AuthServiceAPI{
		authService:     authService,
		loginLimiter:    loginLimiter,
		registerLimiter: registerLimiter,
	}
}

func clientMeta(ctx context.Context) (userAgent, ip string) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ua := md.Get("user-agent"); len(ua) > 0 {
			userAgent = ua[0]
		}
	}
	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		host, _, err := net.SplitHostPort(p.Addr.String())
		if err == nil {
			ip = host
		} else {
			ip = p.Addr.String()
		}
	}
	ip = normalizeRateLimitKey(ip)
	return userAgent, ip
}

func normalizeRateLimitKey(ip string) string {
	if ip == "" {
		return "unknown"
	}

	// Приводим ::1 к 127.0.0.1 (чтобы локальные запросы не плодили разные ключи).
	// Также схлопываем IPv6-mapped IPv4 (::ffff:127.0.0.1) в обычный IPv4.
	if addr, err := netip.ParseAddr(ip); err == nil {
		if addr.IsLoopback() {
			return "127.0.0.1"
		}
		if addr.Is4In6() {
			return addr.Unmap().String()
		}
		if addr.Is4() {
			return addr.String()
		}
	}

	return ip
}

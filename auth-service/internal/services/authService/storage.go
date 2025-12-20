package authService

import (
	"context"
	"time"

	"github.com/vbncursed/medialog/auth-service/internal/models"
)

// Storage — интерфейс хранилища для сервисного слоя auth-service.
// Реализация: `internal/storage/pguserstorage.pguserstorage`.
type Storage interface {
	CreateUser(ctx context.Context, email string, passwordHash string) (uint64, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	CreateSession(ctx context.Context, userID uint64, refreshHash []byte, expiresAt time.Time, userAgent, ip string) (uint64, error)
	GetSessionByRefreshHash(ctx context.Context, refreshHash []byte) (*models.Session, error)
	RevokeSessionByID(ctx context.Context, sessionID uint64, revokedAt time.Time) error
	RevokeAllSessionsByUserID(ctx context.Context, userID uint64, revokedAt time.Time) error
}

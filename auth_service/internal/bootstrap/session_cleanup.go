package bootstrap

import (
	"context"
	"log/slog"
	"time"

	"github.com/vbncursed/medialog/auth-service/internal/storage/auth_storage"
)

// StartSessionCleanup запускает периодическую задачу очистки старых сессий
func StartSessionCleanup(storage *auth_storage.AuthStorage, interval, retentionPeriod time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// Выполняем очистку сразу при старте
		cleanupSessions(storage, retentionPeriod)

		for range ticker.C {
			cleanupSessions(storage, retentionPeriod)
		}
	}()
}

func cleanupSessions(storage *auth_storage.AuthStorage, retentionPeriod time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	deleted, err := storage.CleanupOldSessions(ctx, retentionPeriod)
	if err != nil {
		slog.Error("failed to cleanup old sessions", "err", err)
		return
	}

	if deleted > 0 {
		slog.Info("cleaned up old sessions", "deleted", deleted)
	}
}

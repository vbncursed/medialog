package bootstrap

import (
	"github.com/redis/go-redis/v9"
	"github.com/vbncursed/medialog/auth_service/internal/storage/session_storage"
)

func InitSessionStorage(redisClient *redis.Client) *session_storage.SessionStorage {
	return session_storage.NewSessionStorage(redisClient)
}

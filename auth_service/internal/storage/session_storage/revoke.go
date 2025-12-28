package session_storage

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/vbncursed/medialog/auth_service/internal/models"
)

func (s *SessionStorage) RevokeSessionByRefreshHash(ctx context.Context, refreshHash []byte) error {
	key := sessionKey(refreshHash)
	return s.rdb.Del(ctx, key).Err()
}

func (s *SessionStorage) RevokeAllSessionsByUserID(ctx context.Context, userID uint64) error {
	pattern := sessionKeyPrefix + "*"
	keys, err := s.rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	keysToDelete := make([]string, 0)
	for _, key := range keys {
		data, err := s.rdb.Get(ctx, key).Bytes()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			return err
		}

		var sess models.Session
		if err := json.Unmarshal(data, &sess); err != nil {
			continue
		}

		if sess.UserID == userID {
			keysToDelete = append(keysToDelete, key)
		}
	}

	if len(keysToDelete) > 0 {
		return s.rdb.Del(ctx, keysToDelete...).Err()
	}

	return nil
}

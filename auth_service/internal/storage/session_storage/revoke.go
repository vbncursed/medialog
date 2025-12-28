package session_storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vbncursed/medialog/auth_service/internal/models"
)

func (s *SessionStorage) RevokeSessionByRefreshHash(ctx context.Context, refreshHash []byte) error {
	key := sessionKey(refreshHash)

	sess, err := s.GetSessionByRefreshHash(ctx, refreshHash)
	if err != nil {
		return err
	}

	now := time.Now()
	sess.RevokedAt = &now

	data, err := json.Marshal(sess)
	if err != nil {
		return err
	}

	ttl := time.Until(sess.ExpiresAt)
	if ttl <= 0 {
		ttl = time.Second
	}

	return s.rdb.Set(ctx, key, data, ttl).Err()
}

func (s *SessionStorage) RevokeAllSessionsByUserID(ctx context.Context, userID uint64) error {
	pattern := sessionKeyPrefix + "*"
	keys, err := s.rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	now := time.Now()
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

		if sess.UserID == userID && sess.RevokedAt == nil {
			sess.RevokedAt = &now
			updatedData, err := json.Marshal(sess)
			if err != nil {
				continue
			}

			ttl := time.Until(sess.ExpiresAt)
			if ttl <= 0 {
				ttl = time.Second
			}

			if err := s.rdb.Set(ctx, key, updatedData, ttl).Err(); err != nil {
				return err
			}
		}
	}

	return nil
}

package session_storage

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/vbncursed/medialog/auth/internal/models"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

func (s *SessionStorage) GetSessionByRefreshHash(ctx context.Context, refreshHash []byte) (*models.Session, error) {
	key := sessionKey(refreshHash)
	data, err := s.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	var sess models.Session
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, err
	}

	return &sess, nil
}

package auth_service_api

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter interface {
	Allow(ctx context.Context, key string) bool
}

type redisRateLimiter struct {
	rdb    *redis.Client
	kind   string
	limit  int64
	window time.Duration
}

var incrExpireScript = redis.NewScript(`
local current = redis.call("INCR", KEYS[1])
if current == 1 then
  redis.call("EXPIRE", KEYS[1], ARGV[1])
end
return current
`)

func newRedisRateLimiter(rdb *redis.Client, kind string, limitPerWindow int, window time.Duration) *redisRateLimiter {
	return &redisRateLimiter{
		rdb:    rdb,
		kind:   kind,
		limit:  int64(limitPerWindow),
		window: window,
	}
}

func NewRedisRateLimiter(rdb *redis.Client, kind string, limitPerWindow int, window time.Duration) RateLimiter {
	return newRedisRateLimiter(rdb, kind, limitPerWindow, window)
}

func (l *redisRateLimiter) Allow(ctx context.Context, key string) bool {
	if l == nil || l.limit <= 0 {
		return true
	}
	// Если Redis не инициализирован/недоступен — fail-closed (блокируем auth),
	// потому что без RL нельзя безопасно принимать login/register.
	if l.rdb == nil {
		return false
	}

	if ctx == nil {
		return false
	}
	execCtx := context.WithoutCancel(ctx)

	var cancel context.CancelFunc
	execCtx, cancel = context.WithTimeout(execCtx, 200*time.Millisecond)
	defer cancel()

	redisKey := fmt.Sprintf("rl:%s:%s", l.kind, key)
	ttlSeconds := int64(l.window.Seconds())
	if ttlSeconds <= 0 {
		ttlSeconds = 60
	}

	n, err := incrExpireScript.Run(execCtx, l.rdb, []string{redisKey}, ttlSeconds).Int64()
	if err != nil {
		slog.Warn("rate limit redis error (fail-closed)", "err", err, "key", redisKey)
		return false
	}

	return n <= l.limit
}

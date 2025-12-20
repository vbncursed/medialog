package auth_service_api

import (
	"sync"
	"time"
)

// fixedWindowLimiter — простой rate limiter "N запросов за window" на ключ (например IP).
// Для MVP достаточно, позже можно заменить на Redis-based limiter.
type fixedWindowLimiter struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	state  map[string]fixedWindowState
}

type fixedWindowState struct {
	windowStart time.Time
	count       int
}

func newFixedWindowLimiter(limitPerWindow int, window time.Duration) *fixedWindowLimiter {
	return &fixedWindowLimiter{
		limit:  limitPerWindow,
		window: window,
		state:  make(map[string]fixedWindowState),
	}
}

func (l *fixedWindowLimiter) Allow(key string) bool {
	if l == nil || l.limit <= 0 {
		return true
	}
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	st := l.state[key]
	if st.windowStart.IsZero() || now.Sub(st.windowStart) >= l.window {
		st.windowStart = now
		st.count = 0
	}

	if st.count >= l.limit {
		l.state[key] = st
		return false
	}

	st.count++
	l.state[key] = st
	return true
}

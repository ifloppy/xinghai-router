package app

import (
	"sync"
	"time"
)

const (
	authLoginPerMinute     = 10
	authRegisterPerMinute  = 5
	authEmailCodePerMinute = 5
)

type rateLimiter interface {
	allow(key string) bool
	allowN(key string, n int) bool
	close()
}

type rateWindow struct {
	start time.Time
	count int
}

type memoryLimiter struct {
	mu        sync.Mutex
	perMinute int
	entries   map[string]rateWindow
}

func newMemoryLimiter(n int) *memoryLimiter {
	if n <= 0 {
		n = 60
	}
	return &memoryLimiter{perMinute: n, entries: map[string]rateWindow{}}
}

func newLimiter(n int) *memoryLimiter { return newMemoryLimiter(n) }

func (l *memoryLimiter) allow(key string) bool {
	return l.allowN(key, l.perMinute)
}

func (l *memoryLimiter) allowN(key string, n int) bool {
	if n <= 0 {
		n = l.perMinute
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	w := l.entries[key]
	if now.Sub(w.start) >= time.Minute {
		w = rateWindow{start: now}
	}
	if w.count >= n {
		l.entries[key] = w
		return false
	}
	w.count++
	l.entries[key] = w
	return true
}

func (l *memoryLimiter) close() {}

type fallbackLimiter struct {
	primary *redisLimiter
	backup  *memoryLimiter
}

func (l *fallbackLimiter) allow(key string) bool {
	return l.allowN(key, l.backup.perMinute)
}

func (l *fallbackLimiter) allowN(key string, n int) bool {
	ok, err := l.primary.tryAllowN(key, n)
	if err != nil {
		return l.backup.allowN(key, n)
	}
	return ok
}

func (l *fallbackLimiter) close() {
	l.primary.close()
	l.backup.close()
}

func newRateLimiter(redisURL string, perMinute int) (rateLimiter, string) {
	mem := newMemoryLimiter(perMinute)
	if redisURL == "" {
		return mem, "memory"
	}
	redis, err := newRedisLimiter(redisURL, perMinute)
	if err != nil {
		return mem, "memory"
	}
	return &fallbackLimiter{primary: redis, backup: mem}, "redis"
}

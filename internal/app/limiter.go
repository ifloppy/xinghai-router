package app

import (
	"sync"
	"time"
)

type rateWindow struct {
	start time.Time
	count int
}
type limiter struct {
	mu        sync.Mutex
	perMinute int
	entries   map[string]rateWindow
}

func newLimiter(n int) *limiter { return &limiter{perMinute: n, entries: map[string]rateWindow{}} }
func (l *limiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	w := l.entries[key]
	if now.Sub(w.start) >= time.Minute {
		w = rateWindow{start: now}
	}
	if w.count >= l.perMinute {
		l.entries[key] = w
		return false
	}
	w.count++
	l.entries[key] = w
	return true
}

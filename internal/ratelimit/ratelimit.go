package ratelimit

import (
	"sync"
	"time"
)

type bucket struct {
	minuteCount int
	hourCount   int
	minuteReset time.Time
	hourReset   time.Time
}

// Limiter implements a per-user token bucket rate limiter.
type Limiter struct {
	mu                sync.Mutex
	buckets           map[string]*bucket
	messagesPerMinute int
	messagesPerHour   int
}

// NewLimiter creates a new rate limiter.
func NewLimiter(perMinute, perHour int) *Limiter {
	return &Limiter{
		buckets:           make(map[string]*bucket),
		messagesPerMinute: perMinute,
		messagesPerHour:   perHour,
	}
}

// Allow returns true if the user is allowed to send a message.
func (l *Limiter) Allow(userJID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b, ok := l.buckets[userJID]
	if !ok {
		b = &bucket{
			minuteReset: now.Add(time.Minute),
			hourReset:   now.Add(time.Hour),
		}
		l.buckets[userJID] = b
	}

	if now.After(b.minuteReset) {
		b.minuteCount = 0
		b.minuteReset = now.Add(time.Minute)
	}
	if now.After(b.hourReset) {
		b.hourCount = 0
		b.hourReset = now.Add(time.Hour)
	}

	if b.minuteCount >= l.messagesPerMinute || b.hourCount >= l.messagesPerHour {
		return false
	}

	b.minuteCount++
	b.hourCount++
	return true
}

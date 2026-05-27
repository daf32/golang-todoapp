package core_ratelimit

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type InMemoryLimiter struct {
	rate  rate.Limit
	burst int
	ttl   time.Duration

	mtx     sync.Mutex
	buckets map[string]*bucket
}

type bucket struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func NewMemoryLimiter(
	r rate.Limit,
	burst int,
	ttl time.Duration,
) *InMemoryLimiter {
	l := &InMemoryLimiter{
		rate:    r,
		burst:   burst,
		ttl:     ttl,
		buckets: make(map[string]*bucket),
	}
	go l.cleanUpLoop()
	return l
}

func (l *InMemoryLimiter) Allow(key string) (bool, time.Duration) {
	l.mtx.Lock()
	b, ok := l.buckets[key]
	if !ok {
		b = &bucket{limiter: rate.NewLimiter(l.rate, l.burst)}
		l.buckets[key] = b
	}

	b.lastSeen = time.Now()
	l.mtx.Unlock()

	res := b.limiter.Reserve()
	if !res.OK() {
		return false, 0
	}

	delay := res.Delay()
	if delay > 0 {
		res.Cancel()
		return false, delay
	}

	return true, 0
}

func (l *InMemoryLimiter) cleanUpLoop() {
	ticker := time.NewTicker(l.ttl)
	defer ticker.Stop()
	for range ticker.C {
		cutoff := time.Now().Add(-l.ttl)
		l.mtx.Lock()
		for k, b := range l.buckets {
			if b.lastSeen.Before(cutoff) {
				delete(l.buckets, k)
			}
		}
		l.mtx.Unlock()
	}
}

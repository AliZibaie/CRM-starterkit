package ratelimiter

import (
	"sync"
	"time"
)

type Tracker struct {
	mu          sync.Mutex
	hits        map[string]*bucket
	window      time.Duration
	suspiciousN int
}

type bucket struct {
	count   int
	resetAt time.Time
	blocked bool
}

func New(window time.Duration, suspiciousThreshold int) *Tracker {
	if window <= 0 {
		window = time.Minute
	}
	if suspiciousThreshold <= 0 {
		suspiciousThreshold = 100
	}
	t := &Tracker{hits: make(map[string]*bucket), window: window, suspiciousN: suspiciousThreshold}
	go t.cleanup()
	return t
}

func (t *Tracker) Hit(ip string) (count int, suspicious bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	b, ok := t.hits[ip]
	if !ok || now.After(b.resetAt) {
		t.hits[ip] = &bucket{count: 1, resetAt: now.Add(t.window)}
		return 1, false
	}
	b.count++
	if b.count >= t.suspiciousN {
		b.blocked = true
	}
	return b.count, b.blocked
}

func (t *Tracker) IsSuspicious(ip string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	b, ok := t.hits[ip]
	if !ok {
		return false
	}
	return b.blocked || b.count >= t.suspiciousN
}

func (t *Tracker) MarkBlocked(ip string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := time.Now()
	b, ok := t.hits[ip]
	if !ok || now.After(b.resetAt) {
		t.hits[ip] = &bucket{count: t.suspiciousN, resetAt: now.Add(t.window), blocked: true}
		return
	}
	b.blocked = true
}

func (t *Tracker) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		t.mu.Lock()
		now := time.Now()
		for k, b := range t.hits {
			if now.After(b.resetAt) {
				delete(t.hits, k)
			}
		}
		t.mu.Unlock()
	}
}

package middleware

import (
	"sync"
	"time"

	"crm-starterkit/config"
	"crm-starterkit/internal/core/http/response"

	"github.com/gofiber/fiber/v2"
)

type visitor struct {
	count    int
	resetAt  time.Time
}

type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	max      int
	window   time.Duration
}

func RateLimit(cfg config.RateLimitConfig) fiber.Handler {
	if !cfg.Enabled {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	limiter := &rateLimiter{
		visitors: make(map[string]*visitor),
		max:      cfg.Max,
		window:   cfg.Window,
	}

	go limiter.cleanup()

	return func(c *fiber.Ctx) error {
		key := c.IP()
		if !limiter.allow(key) {
			return response.TooManyRequests(c, "rate limit exceeded")
		}
		return c.Next()
	}
}

func (r *rateLimiter) allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	v, ok := r.visitors[key]
	if !ok || now.After(v.resetAt) {
		r.visitors[key] = &visitor{count: 1, resetAt: now.Add(r.window)}
		return true
	}
	if v.count >= r.max {
		return false
	}
	v.count++
	return true
}

func (r *rateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		r.mu.Lock()
		now := time.Now()
		for key, v := range r.visitors {
			if now.After(v.resetAt) {
				delete(r.visitors, key)
			}
		}
		r.mu.Unlock()
	}
}

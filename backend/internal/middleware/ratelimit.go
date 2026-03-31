package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements a simple per-IP rate limiter
type RateLimiter struct {
	mu       sync.RWMutex
	attempts map[string]*attemptInfo
	limit    int
	window   time.Duration
}

type attemptInfo struct {
	count     int
	windowEnd time.Time
}

// NewRateLimiter creates a new rate limiter
// limit: max attempts per window
// window: time window for rate limiting
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		attempts: make(map[string]*attemptInfo),
		limit:    limit,
		window:   window,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Allow checks if the IP is allowed to make a request
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	info, exists := rl.attempts[ip]
	if !exists || now.After(info.windowEnd) {
		// New window
		rl.attempts[ip] = &attemptInfo{
			count:     1,
			windowEnd: now.Add(rl.window),
		}
		return true
	}

	// Existing window
	if info.count >= rl.limit {
		return false
	}

	info.count++
	return true
}

// cleanup removes expired entries periodically
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, info := range rl.attempts {
			if now.After(info.windowEnd) {
				delete(rl.attempts, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// LoginRateLimit creates a middleware that rate limits login attempts per IP
// Default: 5 attempts per 15 minutes
func LoginRateLimit(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// r.RemoteAddr is already set to the real client IP by Chi's RealIP middleware
			ip := r.RemoteAddr

			if !limiter.Allow(ip) {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", "900") // 15 minutes in seconds
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"status": false, "message": "Too many login attempts. Please try again later."}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

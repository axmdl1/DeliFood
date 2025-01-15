package middleware

import (
	"DeliFood/backend/pkg/logger"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	rps      rate.Limit
	burst    int
	log      *logger.Logger
}

func NewRateLimiter(rps float64, burst int, log *logger.Logger) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      rate.Limit(rps),
		burst:    burst,
		log:      log,
	}
}

// getLimiter retrieves or creates a rate limiter for the given IP
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if limiter, exists := rl.limiters[ip]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(rl.rps, rl.burst)
	rl.limiters[ip] = limiter

	go rl.cleanupLimiter(ip)

	return limiter
}

// cleanupLimiter removes the limiter for the IP after a timeout period
func (rl *RateLimiter) cleanupLimiter(ip string) {
	time.Sleep(5 * time.Minute)
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.limiters, ip)
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			retryAfter := limiter.Reserve().Delay()
			rl.log.Warn("Rate limit exceeded", map[string]interface{}{
				"ip":         ip,
				"url":        r.URL.Path,
				"method":     r.Method,
				"retryAfter": retryAfter.String(),
			})

			w.Header().Set("Retry-After", retryAfter.String())
			http.Error(w, "Rate limit exceeded. Please retry later.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

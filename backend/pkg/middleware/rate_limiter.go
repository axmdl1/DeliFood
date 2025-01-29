package middleware

import (
	"DeliFood/backend/pkg/logger"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	rps      rate.Limit
	burst    int
	log      *logger.Logger
}

// NewRateLimiter initializes a rate limiter middleware
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

// cleanupLimiter removes the limiter for an IP after a timeout period
func (rl *RateLimiter) cleanupLimiter(ip string) {
	time.Sleep(5 * time.Minute)
	rl.mu.Lock()
	delete(rl.limiters, ip)
	rl.mu.Unlock()
}

// LimitMiddleware is a Gin middleware for rate limiting
func (rl *RateLimiter) LimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			retryAfter := limiter.Reserve().Delay()

			rl.log.Warn("Rate limit exceeded", map[string]interface{}{
				"ip":         ip,
				"url":        c.Request.URL.Path,
				"method":     c.Request.Method,
				"retryAfter": retryAfter.String(),
			})

			c.Header("Retry-After", retryAfter.String())
			c.JSON(429, gin.H{"error": "Rate limit exceeded. Please retry later."})
			c.Abort()
			return
		}

		c.Next()
	}
}

package middleware

import (
	"DeliFood/backend/pkg/logger"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter *rate.Limiter
	log     *logger.Logger
}

func NewRateLimiter(rps float64, burst int, log *logger.Logger) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(rate.Limit(rps), burst),
		log:     log,
	}
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.limiter.Allow() {
			rl.log.Warn("Rate limit exceeded", map[string]interface{}{
				"ip":     r.RemoteAddr,
				"url":    r.URL.Path,
				"method": r.Method,
			})
			w.Header().Set("Retry-After", time.Now().Add(time.Second).Format(time.RFC1123))
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

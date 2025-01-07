package cmd

import (
	"DeliFood/backend/pkg/logger"
	"DeliFood/backend/pkg/middleware"
	"net/http"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()

	// Log application start
	log.Info("Application started", map[string]interface{}{
		"module": "main",
		"status": "success",
	})

	// Create Rate Limiter middleware
	rateLimiter := middleware.NewRateLimiter(1, 3, log) // 1 request per second, burst of 3

	// Example HTTP handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Wrap the handler with rate limiting
	http.Handle("/", rateLimiter.Limit(handler))

	// Start HTTP server
	log.Info("Starting server on :8080", nil)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Error("Server failed to start", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

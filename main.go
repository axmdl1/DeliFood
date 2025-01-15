package main

import (
	"DeliFood/backend/handlers"
	"DeliFood/backend/pkg/logger"
	"DeliFood/backend/pkg/middleware"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize logger
	logger := logger.NewLogger()

	// Log application start
	logger.Info("Application started", map[string]interface{}{
		"module": "main",
		"status": "success",
	})

	rateLimiter := middleware.NewRateLimiter(2.0, 5, logger)

	//http mux
	mux := http.NewServeMux()

	//working with server side
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./frontend/assets/"))))

	mux.HandleFunc("/", handlers.MainPageHandler)
	mux.HandleFunc("/contact", handlers.ContactUsHandler)
	mux.HandleFunc("/menu", handlers.MenuHandler)

	rateLimitedMux := rateLimiter.Limit(mux)

	// Start HTTP server
	server := &http.Server{
		Addr:    ":9078",
		Handler: rateLimitedMux,
	}

	//Start server in goroutine
	go func() {
		log.Println("Listening on :9078")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %s\n", err)
		}
	}()

	//Gracefully shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	log.Println("Shutting down the server!")

	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	err := server.Shutdown(shutdownCtx)
	if err != nil {
		logger.Error("Could not shut down the server!!!", map[string]interface{}{"Error:": err})
	}

	logger.Warn("Server gracefully stopped", map[string]interface{}{"Module": "main"})
}

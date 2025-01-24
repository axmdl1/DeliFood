package main

import (
	"DeliFood/backend/handlers"
	"DeliFood/backend/pkg/db"
	"DeliFood/backend/pkg/logger"
	"DeliFood/backend/pkg/middleware"
	"DeliFood/backend/pkg/repo"
	"DeliFood/backend/utils"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	utils.InitJWTSecret()

	// Initialize logger
	logger := logger.NewLogger()
	logger.Info("Application started", map[string]interface{}{
		"module": "main",
		"status": "success",
	})

	// Load database config and connect to the database
	cfg := db.LoadConfigFromEnv(logger)
	dbConn, err := db.NewDB(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Initialize repositories
	userRepo := repo.NewUserRepo(dbConn)
	handlers.SetUserRepo(userRepo)

	// Initialize HTTP mux
	mux := http.NewServeMux()

	// Serve static files (CSS, JS, etc.)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./frontend/assets/"))))

	// Main routes
	mux.HandleFunc("/", handlers.MainPageHandler)
	mux.HandleFunc("/menu", handlers.MenuHandler)
	mux.HandleFunc("/contact", handlers.ContactUsHandler)

	// Auth routes
	authMux := http.NewServeMux()
	authMux.HandleFunc("/register", handlers.RegisterHandler)
	authMux.HandleFunc("/verify-email", handlers.VerifyEmailHandler)
	authMux.HandleFunc("/login", handlers.LoginHandler)
	mux.Handle("/auth/", http.StripPrefix("/auth", authMux))

	// Admin routes
	adminMux := http.NewServeMux()
	adminMux.Handle("/panel", middleware.JWTMiddleware(middleware.RoleMiddleware("admin")(http.HandlerFunc(handlers.AdminPanelHandler))))
	mux.Handle("/admin/", http.StripPrefix("/admin", adminMux))

	// Start the HTTP server
	server := &http.Server{
		Addr:    ":9078",
		Handler: mux,
	}

	go func() {
		log.Println("Server is listening on port 9078")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	// Graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	log.Println("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", map[string]interface{}{"error": err})
	}
	logger.Warn("Server stopped gracefully", nil)
}

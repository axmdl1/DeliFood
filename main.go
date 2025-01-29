package main

import (
	"DeliFood/backend/handlers"
	"DeliFood/backend/pkg/db"
	"DeliFood/backend/pkg/logger"
	"DeliFood/backend/pkg/middleware"
	"DeliFood/backend/pkg/repo"
	"DeliFood/backend/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	// Initialize JWT secret
	utils.InitJWTSecret()

	// Initialize logger
	logger := logger.NewLogger()
	logger.Info("Application started", map[string]interface{}{
		"module": "main",
		"status": "success",
	})

	// Load database config and connect
	cfg := db.LoadConfigFromEnv(logger)
	dbConn, err := db.NewDB(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Initialize repositories
	userRepo := repo.NewUserRepo(dbConn)
	handlers.SetUserRepo(userRepo)

	// Initialize Gin router
	r := gin.Default()

	r.SetFuncMap(utils.TmplFuncs)

	// Load HTML templates
	r.LoadHTMLGlob("frontend/*.html")

	// Serve static assets (CSS, JS, Images)
	r.Static("/assets", "./frontend/assets")

	// Initialize Rate Limiter Middleware (e.g., 5 requests/sec, burst 10)
	rateLimiter := middleware.NewRateLimiter(5, 10, logger)
	r.Use(rateLimiter.LimitMiddleware())

	// Public Routes
	r.GET("/", handlers.MainPageHandler)
	r.GET("/menu", handlers.MenuHandler)
	r.POST("/contact", handlers.ContactUsHandler)

	// Authentication Routes
	authRoutes := r.Group("/auth")
	{
		authRoutes.GET("/register", handlers.RegisterHandler)
		authRoutes.POST("/register", handlers.RegisterHandler)
		authRoutes.GET("/verify-email", handlers.VerifyEmailHandler)
		authRoutes.POST("/verify-email", handlers.VerifyEmailHandler)
		authRoutes.GET("/login", handlers.LoginHandler)
		authRoutes.POST("/login", handlers.LoginHandler)
	}

	// Admin Routes (Protected)
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.RoleMiddleware("admin")) // Only allow admins
	{
		adminRoutes.GET("/panel", handlers.AdminPanelHandler)
	}

	// Start HTTP Server with graceful shutdown
	server := &http.Server{
		Addr:    ":9078",
		Handler: r,
	}

	go func() {
		log.Println("Server is running on port 9078")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	// Graceful Shutdown
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

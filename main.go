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
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	utils.InitJWTSecret()

	appLogger := logger.NewLogger()
	appLogger.Info("Application started", map[string]interface{}{
		"module": "main",
		"status": "success",
	})

	// Connect to MongoDB
	dbClient, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	database := dbClient.Database("GoFood")

	// Initialize repositories using MongoDB collections
	userRepo := repo.NewUserRepo(database.Collection("users"))
	cartRepo := repo.NewCartRepo(database.Collection("cart_items"))
	adminRepo := repo.NewAdminRepo(database.Collection("foods"), database.Collection("users"))
	handlers.SetUserRepo(userRepo)
	handlers.SetCartRepo(cartRepo)
	handlers.SetAdminRepo(adminRepo)

	// Initialize Gin router
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.SetFuncMap(utils.TmplFuncs)

	r.LoadHTMLGlob("frontend/*.html")
	r.Static("/assets", "./frontend/assets")

	rateLimiter := middleware.NewRateLimiter(5, 10, appLogger)
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
		authRoutes.POST("/logout", handlers.LogoutHandler)
	}

	// Cart Routes
	cartRoutes := r.Group("/cart")
	cartRoutes.Use(middleware.AuthMiddleware())
	{
		cartRoutes.GET("/items", handlers.GetCartItemsHandler)
		cartRoutes.POST("/add", handlers.AddToCartHandler)
		cartRoutes.POST("/update/:item_id", handlers.UpdateCartItemHandler)
		cartRoutes.POST("/remove/:item_id", handlers.RemoveCartItemHandler)
		cartRoutes.POST("/checkout", handlers.CheckoutHandler)
	}

	// Admin Routes
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware())
	{
		adminRoutes.GET("/panel", handlers.AdminPanelHandler)
		adminRoutes.POST("/user/change-role", handlers.ChangeUserRoleHandler)
		adminRoutes.POST("/panel/food", handlers.AddFoodHandler)
	}

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

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	log.Println("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("Server shutdown error", map[string]interface{}{"error": err})
	}
	appLogger.Warn("Server stopped gracefully", nil)
}

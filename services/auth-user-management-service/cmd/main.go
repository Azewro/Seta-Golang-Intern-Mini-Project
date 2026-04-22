package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"auth-user-management-service/config"
	"auth-user-management-service/internal/handler"
	"auth-user-management-service/internal/repository"
	"auth-user-management-service/internal/usecase"
)

func main() {
	// Load configuration from .env
	config.LoadEnv()

	// Connect and Auto Migrate Database
	config.ConnectDB()

	// -----------------------------
	// Dependency Injection (Manual)
	// -----------------------------
	// 1. Initialize Repository layer
	userRepo := repository.NewUserRepository(config.DB)

	// 2. Initialize Usecase layer
	authUsecase := usecase.NewAuthUsecase(userRepo)

	// 3. Initialize Handler layer
	authHandler := handler.NewAuthHandler(authUsecase)

	// Initialize Gin App (similar to SpringApplication.run)
	r := gin.Default()

	// Health check API
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "UP",
			"service": "auth-user-management-service",
		})
	})

	// -----------------------------
	// API Routes Grouping
	// -----------------------------
	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			// POST /api/v1/auth/register
			auth.POST("/register", authHandler.Register)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is starting on port http://localhost:%s\n", port)

	// Start the server
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

package main

import (
	"log"
	"os"

	"auth-user-management-service/config"
	_ "auth-user-management-service/docs"
	"auth-user-management-service/internal/handler"
	"auth-user-management-service/internal/middleware"
	"auth-user-management-service/internal/repository"
	"auth-user-management-service/internal/usecase"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Auth & User Management Service API
// @version 1.0
// @description JWT authentication, optional email verification, session revocation, and user APIs. Use header Authorization: Bearer plus your access token.
// @termsOfService http://swagger.io/terms/

// @contact.name Seta Golang Intern Project

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT from POST /api/v1/auth/login (prefix with "Bearer " in Swagger UI Authorize dialog, or paste token only depending on client)

func main() {
	// Load configuration from root .env.backend (or ENV_FILE)
	config.LoadEnv()

	// Connect and Auto Migrate Database
	config.ConnectDB()

	// -----------------------------
	// Dependency Injection (Manual)
	// -----------------------------
	// 1. Initialize Repository layer
	userRepo := repository.NewUserRepository(config.DB)

	// 2. Initialize Repository layer for sessions
	sessionRepo := repository.NewSessionRepository(config.DB)

	// 3. Initialize Repository layer for email verification
	verifyRepo := repository.NewEmailVerificationRepository(config.DB)

	// 4. Initialize Usecase layer
	authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo, verifyRepo)

	// 5. Initialize Handler layer
	authHandler := handler.NewAuthHandler(authUsecase)

	// Initialize Gin App (similar to SpringApplication.run)
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
			auth.POST("/register", authHandler.Register)
			auth.GET("/verify-email", authHandler.VerifyEmail)
			auth.POST("/resend-verification", authHandler.ResendVerification)
			auth.POST("/login", authHandler.Login)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthRequired(sessionRepo))
		{
			protected.POST("/auth/logout", authHandler.Logout)
			protected.GET("/users/me", authHandler.Me)
			protected.POST("/users/bulk", authHandler.BulkGetUsers) // Allow any authenticated user to resolve batch of users

			manager := protected.Group("/")
			manager.Use(middleware.ManagerOnly())
			{
				manager.GET("/users", authHandler.ListUsers)
			}
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is starting on port http://localhost:%s\n", port)
	log.Printf("Swagger UI: http://localhost:%s/swagger/index.html\n", port)

	// Start the server
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"auth-user-management-service/config"
)

func main() {
	// Nạp cấu hình từ .env
	config.LoadEnv()

	// Kết nối và Auto Migrate Databa
	config.ConnectDB()

	// Khởi tạo Gin App (tương tự SpringApplication.run)
	r := gin.Default()

	// API Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "UP",
			"service": "auth-user-management-service",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is starting on port http://localhost:%s\n", port)

	// Khởi chạy server
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

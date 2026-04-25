package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"auth-user-management-service/internal/domain"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadEnv() {
	// Priority order:
	// 1) ENV_FILE explicit path
	// 2) project root .env.backend
	// 3) fallback to system environment variables
	envFile := os.Getenv("ENV_FILE")
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			log.Printf("Failed to load env file from ENV_FILE (%s): %v", envFile, err)
		}
		return
	}

	candidates := []string{
		filepath.Join("..", "..", ".env.backend"),
		".env.backend",
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err != nil {
				log.Printf("Failed to load env file (%s): %v", path, err)
			}
			return
		}
	}

	log.Println("No .env.backend file found, using system environment variables")
}

func ConnectDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL successfully!")
	DB = db

	// Auto Migration: GORM will automatically create/update the table schema
	err = db.AutoMigrate(&domain.User{}, &domain.Session{}, &domain.EmailVerificationToken{})
	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}
	fmt.Println("Auto migration completed!")
}

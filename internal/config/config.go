package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL      string
	JWTSecret        string
	AdminUsername    string
	AdminPasswordHash string
	Port             string
}

func Load() (*Config, error) {
	// .envファイルを読み込む（エラーは無視 - ファイルが存在しない場合は環境変数を使用）
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/contactform?sslmode=disable"),
		JWTSecret:         getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		AdminUsername:     getEnv("ADMIN_USERNAME", "admin"),
		AdminPasswordHash: getEnv("ADMIN_PASSWORD_HASH", ""),
		Port:              getEnv("PORT", "8080"),
	}

	if cfg.AdminPasswordHash == "" {
		return nil, fmt.Errorf("ADMIN_PASSWORD_HASH environment variable is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}


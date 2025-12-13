package main

import (
	"log"

	"contactFormAPI/internal/auth"
	"contactFormAPI/internal/config"
	"contactFormAPI/internal/db"
	"contactFormAPI/internal/http"
	"contactFormAPI/internal/repository"

	_ "contactFormAPI/docs"
)

// @title 問い合わせフォームAPI
// @version 1.0
// @description 問い合わせフォーム用のREST API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT認証。形式: "Bearer {token}"
func main() {
	// 設定の読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// データベース接続
	if err := db.Init(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// JWT認証の初期化
	auth.Init(cfg.JWTSecret)

	// リポジトリの初期化
	contactRepo := repository.NewContactRepository()

	// ルーターの設定
	router := http.SetupRouter(cfg, contactRepo)

	// サーバー起動
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}


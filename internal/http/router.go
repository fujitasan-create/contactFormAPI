package http

import (
	"contactFormAPI/internal/auth"
	"contactFormAPI/internal/config"
	"contactFormAPI/internal/http/handlers"
	"contactFormAPI/internal/http/middleware"
	"contactFormAPI/internal/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(cfg *config.Config, contactRepo *repository.ContactRepository) *gin.Engine {
	router := gin.Default()

	// CORS設定
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4321"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Swagger設定
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 公開API
	public := router.Group("")
	{
		healthHandler := handlers.NewHealthHandler()
		public.GET("/health", healthHandler.HealthCheck)

		contactHandler := handlers.NewContactHandler(contactRepo)
		// POST /contact にレート制限を適用
		public.POST("/contact", middleware.RateLimitMiddleware(), contactHandler.CreateContact)
	}

	// 管理API
	admin := router.Group("/admin")
	{
		adminHandler := handlers.NewAdminHandler(contactRepo, cfg)
		admin.POST("/login", adminHandler.Login)

		// JWT認証が必要なエンドポイント
		adminAuth := admin.Group("")
		adminAuth.Use(auth.JWTAuthMiddleware())
		{
			adminAuth.GET("/messages", adminHandler.GetMessages)
		}
	}

	return router
}

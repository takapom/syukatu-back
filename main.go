package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 設定を読み込む
	config := LoadConfig()
	
	// JWT 鍵を初期化
	InitAuth(config)

	// ① DB 接続＆マイグレーション
	db, err := openGormDB(config)
	if err != nil {
		log.Fatalf("DB 接続エラー: %v", err)
	}

	// ② Gin ルーター初期化
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.CORSAllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// 認証不要ルート
	r.POST("/register", registerHandler(db))
	r.POST("/login", loginHandler(db))

	// 認証ミドルウェアの適用
	auth := r.Group("/")
	auth.Use(authMiddleware())

	// CompanyList 用 CRUD
	auth.POST("/company_lists", createCompanyListHandler(db))
	auth.GET("/company_lists", listCompanyListsHandler(db))
	auth.PUT("/company_lists/:id", updateCompanyListHandler(db))
	auth.DELETE("/company_lists/:id", deleteCompanyListHandler(db))

	// インターンシップ用 CRUD
	auth.POST("/internships", createInternshipHandler(db))
	auth.GET("/internships", listInternshipsHandler(db))
	auth.PUT("/internships/:id", updateInternshipHandler(db))
	auth.DELETE("/internships/:id", deleteInternshipHandler(db))

	// 掲示板用 CRUD
	auth.POST("/posts", createPostHandler(db))
	auth.GET("/posts", getPostsHandler(db))
	auth.GET("/posts/:id", getPostHandler(db))
	auth.DELETE("/posts/:id", deletePostHandler(db))
	auth.POST("/posts/:id/like", likePostHandler(db))
	auth.DELETE("/posts/:id/like", unlikePostHandler(db))
	auth.POST("/posts/:id/comments", createCommentHandler(db))


	// サーバ起動
	addr := fmt.Sprintf(":%s", config.Port)
	log.Printf("Server running on %s (Environment: %s)", addr, config.Environment)
	if err := r.Run(addr); err != nil {
		log.Fatalf("サーバ起動エラー: %v", err)
	}
}

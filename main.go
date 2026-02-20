package main

import (
	"log"
	"time"

	"app/config"
	"app/internal/application/link"
	"app/internal/infrastructure/http"
	"app/internal/infrastructure/persistence/postgres"

	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func router(cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.TrustedPlatform = gin.PlatformCloudflare

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.UIURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Printf("warning: failed to open database: %v", err)
		return r
	}

	if err := db.Ping(); err != nil {
		log.Printf("warning: failed to ping database: %v", err)
		return r
	}

	if err := goose.Up(db, "db/migrations"); err != nil {
		log.Printf("warning: failed to run migrations: %v", err)
	}

	repo := postgres.NewLinkRepository(db)
	service := link.NewService(repo, cfg.BaseURL)

	handler := http.NewHandler(service)
	handler.RegisterRoutes(r)

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	return r
}

func main() {
	cfg := config.Load()

	r := router(cfg)

	if err := r.Run(":" + cfg.Port); err != nil {
		panic(err)
	}
}

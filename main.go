package main

import (
	"log"
	"os"
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

func router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.TrustedPlatform = gin.PlatformCloudflare

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	cfg := config.Load()

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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router().Run(":" + port); err != nil {
		panic(err)
	}
}

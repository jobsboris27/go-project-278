package main

import (
	"database/sql"
	"log"
	"os"

	"app/config"
	"app/internal/application/link"
	"app/internal/infrastructure/http"
	"app/internal/infrastructure/persistence/postgres"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	if err := goose.Up(db, "db/migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	repo := postgres.NewLinkRepository(db)
	service := link.NewService(repo, cfg.BaseURL)

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	handler := http.NewHandler(service)
	handler.RegisterRoutes(router)

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

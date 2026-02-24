package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"app/config"
	"app/internal/application/link"
	"app/internal/infrastructure/http"
	"app/internal/infrastructure/persistence/postgres"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"

	"github.com/rollbar/rollbar-go"
)

func doSomething() {
	var timer *time.Timer = nil
	timer.Reset(10) // this will panic
}

func connectDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func createDependencies(db *sql.DB, baseURL string) *link.Service {
	repo := postgres.NewLinkRepository(db)
	return link.NewService(repo, baseURL)
}

func initRollbar(token string) {
	rollbar.SetToken(token)
	rollbar.SetEnvironment("production")
}

func registerRoutes(r *gin.Engine, service *link.Service) {
	handler := http.NewHandler(service)
	handler.RegisterRoutes(r)

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}

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

	return r
}

func main() {
	cfg := config.Load()

	initRollbar(cfg.RollbarToken)
	defer rollbar.Close()

	rollbar.Info("Application starting")
	rollbar.WrapAndWait(doSomething)

	db, err := connectDB(cfg.DatabaseURL)
	if err != nil {
		log.Printf("error: failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error: failed to close database: %v", err)
		}
	}()

	service := createDependencies(db, cfg.BaseURL)

	r := router(cfg)
	registerRoutes(r, service)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Printf("error: failed to start server: %v", err)
		os.Exit(1)
	}
}

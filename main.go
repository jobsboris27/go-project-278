// Package main provides the main entry point for the application.
package main

import (
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func router() *gin.Engine {
	r := gin.New()

	if dsn := os.Getenv("SENTRY_DSN"); dsn != "" {
		if err := sentry.Init(sentry.ClientOptions{Dsn: dsn}); err != nil {
			panic(err)
		}
		r.Use(sentrygin.New(sentrygin.Options{}))
		defer sentry.Flush(2 * time.Second)
	}

	r.Use(gin.Logger(), gin.Recovery())
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("/panic", func(_ *gin.Context) {
		panic("test error for Sentry")
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

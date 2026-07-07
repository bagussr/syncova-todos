package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syncova-todo/cmd/docs"
	"syncova-todo/config"
	"syncova-todo/middleware"
	"syscall"
	"time"

	delivery "syncova-todo/delivery/http/v1"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	docs.SwaggerInfo.Title = "Syncova Todo API"
	docs.SwaggerInfo.Description = "This is a sample server for Syncova Todo API."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = cfg.Host + ":" + cfg.Port
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	app := gin.Default()
	app.Use(middleware.ErrorMiddleware())
	app.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
	))
	api := app.Group("/api")
	delivery.SetupRouter(api)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: app,
	}

	// Start server in goroutine
	go func() {
		log.Println("🚀 Server running on", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("⚠️ Server forced to shutdown:", err)
	} else {
		log.Println("✅ Server exited gracefully")
	}

	os.Exit(0)
}

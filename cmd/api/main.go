package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syncova-todo/config"
	delivery "syncova-todo/delivery/http/v1"
	"syncova-todo/middleware"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	app := gin.Default()
	app.Use(middleware.ErrorMiddleware())
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

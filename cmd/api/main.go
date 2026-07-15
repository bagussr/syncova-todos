package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syncova-todo/cmd/docs"
	"syncova-todo/config"
	domain "syncova-todo/domain/base"
	clients "syncova-todo/infrastructure/clients/auth"
	"syncova-todo/infrastructure/database"
	"syncova-todo/middleware"
	"syscall"
	"time"

	delivery "syncova-todo/delivery/http/v1"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your JWT token.

func main() {
	cfg := config.LoadConfig()

	// configuration swagger
	docs.SwaggerInfo.Title = "Syncova Todo API"
	docs.SwaggerInfo.Description = "This is a sample server for Syncova Todo API."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = cfg.Host + ":" + cfg.Port
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// initialize the database connection
	database, err := database.NewPostgresConnection(cfg, true)

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	defer func() {
		sqlDB, err := database.DB.DB()
		if err != nil {
			log.Printf("Failed to get sqlDB: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
	}()

	// calling the auth middleware
	authClient := clients.NewAuthClientWithAPIKey(
		cfg.AuthServiceURL,
		5*time.Second,
		cfg.AuthServiceAPIKey,
	)

	app := gin.Default()
	app.Use(middleware.ErrorMiddleware())
	app.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
	))

	app.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, domain.BaseResponse{
			Success:    true,
			StatusCode: 200,
			Message:    "Server is running",
		})
	})

	api := app.Group("/api")

	api.Use(middleware.RequireAuth(authClient))
	delivery.SetupRouter(api, database)

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

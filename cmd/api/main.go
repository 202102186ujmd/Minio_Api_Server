package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"github.com/202102186ujmd/Minio_Api_Server/internal/config"
	"github.com/202102186ujmd/Minio_Api_Server/internal/domain/storage"
	"github.com/202102186ujmd/Minio_Api_Server/internal/http/handlers"
	"github.com/202102186ujmd/Minio_Api_Server/internal/http/middleware"
	"github.com/202102186ujmd/Minio_Api_Server/internal/http/routes"
	"github.com/202102186ujmd/Minio_Api_Server/internal/infrastructure/minio"
	"github.com/202102186ujmd/Minio_Api_Server/internal/utils/logger"
	_ "github.com/202102186ujmd/Minio_Api_Server/docs"
)

// @title MinIO Management API
// @version 1.1.0
// @description API profesional para gestión de MinIO (CRUD de objetos), con validaciones, logging, métricas y respuestas estandarizadas.
// @termsOfService https://example.com/terms
// @contact.name Soporte API
// @contact.url https://example.com/support
// @contact.email soporte@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	logg, err := logger.New(cfg)
	if err != nil {
		log.Fatalf("logger error: %v", err)
	}
	defer func() {
		_ = logg.Sync()
	}()

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	minioClient, err := minio.NewClient(cfg)
	if err != nil {
		logg.Fatal("minio init failed", logger.Error(err))
	}

	storageService := storage.NewService(minioClient, cfg)

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(logg))
	r.Use(middleware.Recovery(logg))
	r.Use(middleware.CORS(cfg))
	r.Use(middleware.Metrics())
	r.Use(middleware.APIKey(cfg))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.Register(r,
		handlers.NewHealthHandler(storageService, cfg),
		handlers.NewStorageHandler(storageService, cfg),
	)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		logg.Info("server starting", logger.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Fatal("server error", logger.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logg.Info("server shutting down")
	if err := srv.Shutdown(ctx); err != nil {
		logg.Error("server shutdown error", logger.Error(err))
	}
}

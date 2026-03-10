package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bambu-farm/api"
	"bambu-farm/pkg/config"
	"bambu-farm/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize configuration
	cfg := config.LoadConfig()

	// Initialize logger
	log := logger.InitLogger(cfg.Env)
	defer log.Sync()

	log.Info("Starting BambuLab Print Farm Manager Backend...")

	// Set Gin mode
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Database
	db := repository.InitDB()

	// Initialize Repositories
	authRepo := repository.NewAuthRepository(db)
	printerRepo := repository.NewPrinterRepository(db)

	// Initialize Services
	authService := service.NewAuthService(authRepo)
	printerService := service.NewPrinterService(printerRepo)

	// Initialize Handlers
	authHandler := api.NewAuthHandler(authService)
	printerHandler := api.NewPrinterHandler(printerService)

	// Initialize router
	router := gin.Default()

	// Register generic routes
	api.RegisterRoutes(router)

	// Register Feature routes
	authHandler.RegisterRoutes(router)
	printerHandler.RegisterRoutes(router)

	// Setup Server
	addr := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Run Server in goroutine
	go func() {
		log.Infof("Server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %s", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exiting")
}

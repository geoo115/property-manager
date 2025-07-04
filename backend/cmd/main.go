package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/geoo115/property-manager/config"
	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/events"
	"github.com/geoo115/property-manager/logger"
	"github.com/geoo115/property-manager/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title Property Management System API
// @version 1.0
// @description A comprehensive property management system with role-based access control
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger.InitLogger()
	logger.LogInfo("Starting Property Management System", logrus.Fields{
		"version": "1.0.0",
		"env":     cfg.Server.GinMode,
	})

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Initialize database (includes Redis initialization)
	if err := db.Init(cfg); err != nil {
		logger.LogError(err, "Failed to initialize database", nil)
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Initialize Kafka
	if err := events.InitKafka(cfg); err != nil {
		logger.LogError(err, "Failed to initialize Kafka", nil)
		log.Fatalf("Kafka initialization failed: %v", err)
	}

	// Start Kafka Consumer in a separate goroutine
	logger.LogInfo("Starting Kafka Consumer", nil)
	go func() {
		if err := events.StartKafkaConsumer(cfg); err != nil {
			logger.LogError(err, "Kafka consumer failed", nil)
		}
	}()

	// Setup Gin router
	r := gin.New()

	// Add middleware
	r.Use(logger.GinLogger())
	r.Use(logger.GinRecovery())

	// Configure CORS
	corsConfig := cors.Config{
		AllowOrigins:     cfg.Security.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	// Setup API routes
	router.SetupRouter(r, cfg)

	// Setup server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.LogInfo("Server starting", logrus.Fields{
			"address": serverAddr,
			"mode":    cfg.Server.GinMode,
		})

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.LogError(err, "Server failed to start", nil)
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.LogInfo("Server shutting down...", nil)

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		logger.LogError(err, "Server forced to shutdown", nil)
		log.Fatal("Server forced to shutdown:", err)
	}

	// Close database connection
	if err := db.Close(); err != nil {
		logger.LogError(err, "Failed to close database connection", nil)
	}

	// Close Kafka producer
	events.CloseKafka()

	logger.LogInfo("Server exited gracefully", nil)
}

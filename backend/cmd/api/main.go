package main

import (
	"log"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/logger"

	"github.com/gin-gonic/gin"
)

// @title CruiseBooking API
// @version 1.0.0
// @description 邮轮舱位预订平台 RESTful API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@cruisebooking.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	log := logger.New(cfg.LogLevel)
	defer log.Sync()

	// Initialize database
	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", log.Error(err))
	}

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	r := gin.New()

	// TODO: Setup middleware
	// TODO: Setup routes
	// TODO: Start server

	log.Info("Server starting on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server", log.Error(err))
	}
}

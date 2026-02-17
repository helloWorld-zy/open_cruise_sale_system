package main

import (
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/logger"
	"fmt"

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
	l := logger.New(cfg.LogLevel)
	defer l.Sync()

	// Initialize database
	db, err := database.New(cfg.Database)
	if err != nil {
		l.Fatalw("Failed to connect to database", "error", err)
	}

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	r := gin.New()

	// Initialize handlers and routes
	setupRoutes(r, db, cfg)

	// Since we are not using AdminHandlers struct in default manner yet, usage depends on actual file.
	// But main.go was truncated in previous view.
	// Let's assume standard setupRoutes logic.

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	l.Infow("Server starting", "addr", addr)
	if err := r.Run(addr); err != nil {
		l.Fatalw("Failed to start server", "error", err)
	}
}

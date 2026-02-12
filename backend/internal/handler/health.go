package handler

import (
	"backend/internal/database"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db *database.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// RegisterRoutes registers health check routes
func (h *HealthHandler) RegisterRoutes(router *gin.RouterGroup) {
	health := router.Group("/health")
	{
		health.GET("", h.BasicHealthCheck)
		health.GET("/ready", h.ReadinessCheck)
		health.GET("/live", h.LivenessCheck)
	}
}

// BasicHealthCheck performs a basic health check
func (h *HealthHandler) BasicHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "cruise-booking-api",
		"version":   "1.0.0",
	})
}

// ReadinessCheck checks if the service is ready to accept traffic
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	status := gin.H{
		"status":    "ready",
		"timestamp": time.Now().UTC(),
		"checks":    make(map[string]interface{}),
	}

	checks := status["checks"].(map[string]interface{})
	allHealthy := true

	// Check database connection
	dbHealthy := h.checkDatabase()
	checks["database"] = gin.H{
		"status":  map[bool]string{true: "healthy", false: "unhealthy"}[dbHealthy],
		"message": map[bool]string{true: "Connected", false: "Disconnected"}[dbHealthy],
	}
	if !dbHealthy {
		allHealthy = false
	}

	if !allHealthy {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "not_ready",
			"timestamp": time.Now().UTC(),
			"checks":    checks,
		})
		return
	}

	c.JSON(http.StatusOK, status)
}

// LivenessCheck checks if the service is alive
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now().UTC(),
	})
}

// checkDatabase verifies database connection
func (h *HealthHandler) checkDatabase() bool {
	if h.db == nil {
		return false
	}

	sqlDB, err := h.db.DB.DB()
	if err != nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx) == nil
}

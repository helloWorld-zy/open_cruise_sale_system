package main

import (
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// setupAdminRoutes configures admin API routes with RBAC protection
func setupAdminRoutes(r *gin.Engine, handlers *AdminHandlers, cfg *config.Config) {
	// Admin API group with authentication and authorization
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.JWTAuth(&cfg.JWT))
	admin.Use(middleware.RequireRole("super_admin", "operations"))
	{
		// Cruise management
		cruises := admin.Group("/cruises")
		{
			cruises.GET("", handlers.AdminCruise.List)
			cruises.POST("", handlers.AdminCruise.Create)
			cruises.GET("/:id", handlers.AdminCruise.GetByID)
			cruises.PUT("/:id", handlers.AdminCruise.Update)
			cruises.DELETE("/:id", handlers.AdminCruise.Delete)
			cruises.POST("/:id/restore", handlers.AdminCruise.Restore)
			cruises.PUT("/:id/status", handlers.AdminCruise.UpdateStatus)
			cruises.POST("/upload", handlers.AdminCruise.UploadImages)
		}

		// Cabin type management
		cabinTypes := admin.Group("/cabin-types")
		{
			cabinTypes.GET("", handlers.AdminCabinType.List)
			cabinTypes.POST("", handlers.AdminCabinType.Create)
			cabinTypes.GET("/:id", handlers.AdminCabinType.GetByID)
			cabinTypes.PUT("/:id", handlers.AdminCabinType.Update)
			cabinTypes.DELETE("/:id", handlers.AdminCabinType.Delete)
			cabinTypes.POST("/:id/restore", handlers.AdminCabinType.Restore)
			cabinTypes.PUT("/:id/status", handlers.AdminCabinType.UpdateStatus)
		}

		// Facility management
		facilities := admin.Group("/facilities")
		{
			facilities.GET("", handlers.AdminFacility.List)
			facilities.POST("", handlers.AdminFacility.Create)
			facilities.GET("/:id", handlers.AdminFacility.GetByID)
			facilities.PUT("/:id", handlers.AdminFacility.Update)
			facilities.DELETE("/:id", handlers.AdminFacility.Delete)
		}

		// Facility category management
		categories := admin.Group("/facility-categories")
		{
			categories.GET("", handlers.AdminFacilityCategory.ListCategories)
			categories.POST("", handlers.AdminFacilityCategory.CreateCategory)
			categories.PUT("/:id", handlers.AdminFacilityCategory.UpdateCategory)
			categories.DELETE("/:id", handlers.AdminFacilityCategory.DeleteCategory)
		}
	}
}

// AdminHandlers groups all admin handlers
type AdminHandlers struct {
	AdminCruise           *handler.AdminCruiseHandler
	AdminCabinType        *handler.AdminCabinTypeHandler
	AdminFacility         *handler.AdminFacilityHandler
	AdminFacilityCategory *handler.AdminFacilityCategoryHandler
}

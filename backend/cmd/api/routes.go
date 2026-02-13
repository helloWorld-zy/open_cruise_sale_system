package main

import (
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/repository"
	"backend/internal/service"
	"backend/internal/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// setupRoutes configures all API routes
func setupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Initialize repositories
	cruiseRepo := repository.NewCruiseRepository(db)
	cabinTypeRepo := repository.NewCabinTypeRepository(db)
	facilityRepo := repository.NewFacilityRepository(db)
	facilityCategoryRepo := repository.NewFacilityCategoryRepository(db)

	// Initialize services
	cruiseService := service.NewCruiseService(cruiseRepo)
	cabinTypeService := service.NewCabinTypeService(cabinTypeRepo)
	facilityService := service.NewFacilityService(facilityRepo, facilityCategoryRepo)
	facilityCategoryService := service.NewFacilityCategoryService(facilityCategoryRepo)

	// Initialize MinIO client
	minioClient, err := storage.New(cfg.MinIO)
	if err != nil {
		// Log error but continue (or panic depending on requirement)
		// For now we panic to ensure valid configuration
		panic(err)
	}
	storageService := service.NewStorageService(minioClient, cfg.MinIO.Endpoint)

	// Initialize handlers
	cruiseHandler := handler.NewCruiseHandler(cruiseService)
	cabinTypeHandler := handler.NewCabinTypeHandler(cabinTypeService)
	facilityHandler := handler.NewFacilityHandler(facilityService)
	facilityCategoryHandler := handler.NewFacilityCategoryHandler(facilityCategoryService)

	// Initialize admin handlers
	adminHandlers := &AdminHandlers{
		AdminCruise:           handler.NewAdminCruiseHandler(cruiseService, storageService),
		AdminCabinType:        handler.NewAdminCabinTypeHandler(cabinTypeService),
		AdminFacility:         handler.NewAdminFacilityHandler(facilityService),
		AdminFacilityCategory: handler.NewAdminFacilityCategoryHandler(facilityCategoryService),
	}

	// Setup admin routes
	setupAdminRoutes(r, adminHandlers, cfg)

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Cruise routes
		cruises := v1.Group("/cruises")
		{
			cruises.GET("", cruiseHandler.List)
			cruises.POST("", cruiseHandler.Create)
			cruises.GET("/code/:code", cruiseHandler.GetByCode)
			cruises.GET("/:id", cruiseHandler.GetByID)
			cruises.PUT("/:id", cruiseHandler.Update)
			cruises.DELETE("/:id", cruiseHandler.Delete)
			cruises.POST("/:id/restore", cruiseHandler.Restore)
			cruises.PUT("/:id/status", cruiseHandler.UpdateStatus)
			cruises.PUT("/:id/sort-weight", cruiseHandler.UpdateSortWeight)

			// Cabin types nested under cruise
			cruises.GET("/:cruise_id/cabin-types", cabinTypeHandler.ListByCruise)
			cruises.GET("/:cruise_id/facilities", facilityHandler.ListByCruise)
			cruises.GET("/:cruise_id/facilities/grouped", facilityHandler.ListByCruiseGrouped)
			cruises.GET("/:cruise_id/facility-categories", facilityCategoryHandler.ListCategoriesByCruise)
		}

		// Cabin type routes
		cabinTypes := v1.Group("/cabin-types")
		{
			cabinTypes.GET("", cabinTypeHandler.List)
			cabinTypes.POST("", cabinTypeHandler.Create)
			cabinTypes.GET("/:id", cabinTypeHandler.GetByID)
			cabinTypes.PUT("/:id", cabinTypeHandler.Update)
			cabinTypes.DELETE("/:id", cabinTypeHandler.Delete)
			cabinTypes.POST("/:id/restore", cabinTypeHandler.Restore)
			cabinTypes.PUT("/:id/status", cabinTypeHandler.UpdateStatus)
		}

		// Facility routes
		facilities := v1.Group("/facilities")
		{
			facilities.GET("", facilityHandler.List)
			facilities.POST("", facilityHandler.Create)
			facilities.GET("/:id", facilityHandler.GetByID)
			facilities.PUT("/:id", facilityHandler.Update)
			facilities.DELETE("/:id", facilityHandler.Delete)
		}

		// Facility category routes
		facilityCategories := v1.Group("/facility-categories")
		{
			facilityCategories.GET("", facilityCategoryHandler.ListCategoriesByCruise)
			facilityCategories.POST("", facilityCategoryHandler.CreateCategory)
			facilityCategories.PUT("/:id", facilityCategoryHandler.UpdateCategory)
			facilityCategories.DELETE("/:id", facilityCategoryHandler.DeleteCategory)
		}
	}
}

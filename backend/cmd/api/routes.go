package main

import (
	"backend/internal/auth"
	"backend/internal/cache"
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/messaging"
	"backend/internal/middleware"
	"backend/internal/payment"
	"backend/internal/repository"
	"backend/internal/service"
	"backend/internal/storage"
	"os"

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
	voyageRepo := repository.NewVoyageRepository(db)
	cabinRepo := repository.NewCabinRepository(db)
	priceRepo := repository.NewPriceRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Initialize infrastructure clients (best-effort)
	var redisClient *cache.RedisClient
	if client, err := cache.New(cfg.Redis); err == nil {
		redisClient = client
	}

	var natsConn *messaging.NATSClient
	if client, err := messaging.New(cfg.NATS); err == nil {
		natsConn = client
	}

	// Initialize services
	cruiseService := service.NewCruiseService(cruiseRepo)
	cabinTypeService := service.NewCabinTypeService(cabinTypeRepo)
	facilityService := service.NewFacilityService(facilityRepo, facilityCategoryRepo)
	facilityCategoryService := service.NewFacilityCategoryService(facilityCategoryRepo)

	jwtService := service.NewJWTService(&cfg.JWT)
	wechatAuthService := service.NewWechatAuthService(service.WechatAuthConfig{
		AppID:     cfg.Wechat.AppID,
		AppSecret: os.Getenv("WECHAT_APP_SECRET"),
	}, userRepo, jwtService)
	smsService := service.NewSMSService(service.SMSConfig{Provider: "mock"}, userRepo)

	var tokenBlacklist middleware.TokenBlacklist
	if redisClient != nil {
		tokenBlacklist = middleware.NewTokenBlacklist(redisClient.GetClient())
	}

	rbac, _ := auth.NewRBAC()

	orderService := func() service.OrderService {
		if redisClient != nil {
			return service.NewOrderService(orderRepo, voyageRepo, cabinRepo, priceRepo, inventoryRepo, redisClient.GetClient())
		}
		return service.NewOrderService(orderRepo, voyageRepo, cabinRepo, priceRepo, inventoryRepo)
	}()

	paymentService := func() payment.PaymentService {
		if natsConn != nil && redisClient != nil {
			return payment.NewPaymentService(orderRepo, natsConn.GetConn(), redisClient.GetClient())
		}
		if natsConn != nil {
			return payment.NewPaymentService(orderRepo, natsConn.GetConn())
		}
		if redisClient != nil {
			return payment.NewPaymentService(orderRepo, nil, redisClient.GetClient())
		}
		return payment.NewPaymentService(orderRepo, nil)
	}()

	wechatConfig := payment.WechatPayConfig{
		AppID:    cfg.Wechat.AppID,
		MchID:    cfg.Wechat.MchID,
		APIv3Key: cfg.Wechat.APIV3Key,
		NotifyURL: func() string {
			if v := os.Getenv("WECHAT_NOTIFY_URL"); v != "" {
				return v
			}
			return "http://localhost:8080/api/v1/payments/callback/wechat"
		}(),
		Certificate: func() string {
			if cfg.Wechat.CertPath == "" {
				return ""
			}
			b, err := os.ReadFile(cfg.Wechat.CertPath)
			if err != nil {
				return ""
			}
			return string(b)
		}(),
		PrivateKey: func() string {
			if cfg.Wechat.KeyPath == "" {
				return ""
			}
			b, err := os.ReadFile(cfg.Wechat.KeyPath)
			if err != nil {
				return ""
			}
			return string(b)
		}(),
		SerialNo: os.Getenv("WECHAT_SERIAL_NO"),
	}

	var wechatProvider payment.PaymentProvider
	if redisClient != nil {
		wechatProvider = payment.NewWechatPay(wechatConfig, orderRepo, redisClient.GetClient())
	} else {
		wechatProvider = payment.NewWechatPay(wechatConfig, orderRepo)
	}

	paymentService.RegisterProvider("wechat", wechatProvider)

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
	authHandler := handler.NewAuthHandler(&cfg.JWT, userRepo, rbac, tokenBlacklist)
	userHandler := handler.NewUserHandler(wechatAuthService, smsService, userRepo)
	orderHandler := handler.NewOrderHandler(orderService)
	orderQueryHandler := handler.NewOrderQueryHandler(orderService, orderRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService)

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

		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.Refresh)
			authGroup.POST("/wechat-login", authHandler.WeChatLogin)
			authGroup.POST("/sms/send", userHandler.SendSMSCode)
			authGroup.POST("/sms/login", userHandler.SMSLogin)
			authGroup.POST("/wechat/login", userHandler.WechatLogin)
			authGroup.POST("/wechat/phone-login", userHandler.WechatPhoneLogin)

			authProtected := authGroup.Group("")
			authProtected.Use(middleware.JWTAuth(&cfg.JWT))
			{
				authProtected.POST("/logout", authHandler.Logout)
				authProtected.GET("/me", authHandler.Me)
				authProtected.POST("/change-password", authHandler.ChangePassword)
			}
		}

		orders := v1.Group("/orders")
		orders.Use(middleware.JWTAuth(&cfg.JWT))
		{
			orders.POST("", orderHandler.Create)
			orders.POST("/calculate", orderHandler.Calculate)
			orders.GET("", orderHandler.List)
			orders.GET("/my", orderQueryHandler.GetMyOrders)
			orders.GET("/statistics", orderQueryHandler.GetOrderStatistics)
			orders.GET("/:id", orderHandler.GetByID)
			orders.GET("/:id/detail", orderQueryHandler.GetOrderDetail)
			orders.GET("/:id/details", orderHandler.GetWithDetails)
			orders.GET("/number/:orderNumber", orderHandler.GetByOrderNumber)
			orders.PUT("/:id", orderHandler.Update)
			orders.POST("/:id/cancel", orderHandler.Cancel)
			orders.POST("/:id/confirm", orderHandler.Confirm)
			orders.POST("/:id/complete", orderHandler.Complete)
			orders.DELETE("/:id", orderHandler.Delete)
		}

		payments := v1.Group("/payments")
		{
			payments.POST("/callback/wechat", paymentHandler.WechatCallback)

			paymentsProtected := payments.Group("")
			paymentsProtected.Use(middleware.JWTAuth(&cfg.JWT))
			{
				paymentsProtected.POST("", paymentHandler.Create)
				paymentsProtected.GET("/:id", paymentHandler.Query)
				paymentsProtected.GET("/order/:orderId", paymentHandler.GetByOrder)
				paymentsProtected.POST("/:id/refund", paymentHandler.Refund)
			}
		}

		user := v1.Group("/user")
		user.Use(middleware.JWTAuth(&cfg.JWT))
		{
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
			user.GET("/passengers", userHandler.ListFrequentPassengers)
			user.POST("/passengers", userHandler.CreateFrequentPassenger)
			user.PUT("/passengers/:id", userHandler.UpdateFrequentPassenger)
			user.DELETE("/passengers/:id", userHandler.DeleteFrequentPassenger)
		}

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

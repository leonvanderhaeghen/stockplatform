package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "github.com/leonvanderhaeghen/stockplatform/services/gatewaySvc/docs" // Import generated docs
	"github.com/leonvanderhaeghen/stockplatform/services/gatewaySvc/internal/services"
)

// Server represents the REST API server
type Server struct {
	router      *gin.Engine
	productSvc  services.ProductService
	inventorySvc services.InventoryService
	orderSvc    services.OrderService
	userSvc     services.UserService
	supplierSvc services.SupplierService
	logger      *zap.Logger
	jwtSecret   string
	port        string
}

// NewServer creates a new REST API server
func NewServer(
	productSvc services.ProductService,
	inventorySvc services.InventoryService,
	orderSvc services.OrderService,
	userSvc services.UserService,
	supplierSvc services.SupplierService,
	jwtSecret string,
	port string,
	logger *zap.Logger,
) *Server {
	router := gin.New()
	
	// Add middlewares
	router.Use(gin.Recovery())
	router.Use(loggerMiddleware(logger))
	
	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	
	return &Server{
		router:      router,
		productSvc:  productSvc,
		inventorySvc: inventorySvc,
		orderSvc:    orderSvc,
		userSvc:     userSvc,
		supplierSvc: supplierSvc,
		logger:      logger.Named("rest_server"),
		jwtSecret:   jwtSecret,
		port:        port,
	}
}

// SetupRoutes configures all API routes
func (s *Server) SetupRoutes() {
	// API versioning
	v1 := s.router.Group("/api/v1")
	
	// Health check
	v1.GET("/health", s.healthCheck)
	
	// Swagger documentation
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, 
		ginSwagger.URL("/swagger/doc.json"), // The url pointing to API definition
		ginSwagger.DefaultModelsExpandDepth(-1), // Hide models section
	))
	
	// Authentication routes
	auth := v1.Group("/auth")
	{
		auth.POST("/register", s.registerUser)
		auth.POST("/login", s.loginUser)
	}
	
	// User routes (protected)
	users := v1.Group("/users")
	users.Use(s.authMiddleware())
	{
		// Public user routes (for admin)
		users.GET("", s.listUsers)
		users.GET("/:id", s.getUserByID)
		
		// Current user routes
		users.GET("/me", s.getCurrentUser)
		users.PUT("/me", s.updateUserProfile)
		users.PUT("/me/password", s.changeUserPassword)
		
		// Address management
		users.GET("/me/addresses", s.getUserAddresses)
		users.POST("/me/addresses", s.createUserAddress)
		users.GET("/me/addresses/default", s.getUserDefaultAddress)
		users.PUT("/me/addresses/:id", s.updateUserAddress)
		users.DELETE("/me/addresses/:id", s.deleteUserAddress)
		users.PUT("/me/addresses/:id/default", s.setDefaultUserAddress)
	}
	
	// Admin routes (protected + admin role)
	admin := v1.Group("/admin")
	admin.Use(s.authMiddleware(), s.adminMiddleware())
	{
		admin.GET("/users", s.listUsers)
		admin.GET("/users/:id", s.getUserByID)
		admin.PUT("/users/:id/activate", s.activateUser)
		admin.PUT("/users/:id/deactivate", s.deactivateUser)
	}
	
	// Product routes
	products := v1.Group("/products")
	{
		products.GET("", s.listProducts)
		products.GET("/:id", s.getProduct)
		products.GET("/categories", s.listCategories)
		
		// Protected product routes (admin/staff only)
		productsAdmin := products.Group("")
		productsAdmin.Use(s.authMiddleware(), s.staffMiddleware())
		{
			productsAdmin.POST("", s.createProduct)
			productsAdmin.PUT("/:id", s.updateProduct)
			productsAdmin.DELETE("/:id", s.deleteProduct)
			productsAdmin.POST("/categories", s.createCategory)
		}
	}
	
	// Inventory routes (mostly protected)
	inventory := v1.Group("/inventory")
	inventory.Use(s.authMiddleware(), s.staffMiddleware())
	{
		inventory.GET("", s.listInventory)
		inventory.GET("/:id", s.getInventoryItem)
		inventory.GET("/product/:productId", s.getInventoryItemByProduct)
		inventory.GET("/sku/:sku", s.getInventoryItemBySKU)
		inventory.POST("", s.createInventoryItem)
		inventory.PUT("/:id", s.updateInventoryItem)
		inventory.DELETE("/:id", s.deleteInventoryItem)
		inventory.POST("/:id/stock/add", s.addStock)
		inventory.POST("/:id/stock/remove", s.removeStock)
	}
	
	// Order routes
	orders := v1.Group("/orders")
	orders.Use(s.authMiddleware())
	{
		// Customer routes
		orders.GET("/me", s.getUserOrders)
		orders.GET("/me/:id", s.getUserOrder)
		orders.POST("", s.createOrder)
		
		// Admin/staff routes
		ordersAdmin := orders.Group("")
		ordersAdmin.Use(s.staffMiddleware())
		{
			ordersAdmin.GET("", s.listOrders)
			ordersAdmin.GET("/:id", s.getOrder)
			ordersAdmin.PUT("/:id/status", s.updateOrderStatus)
			ordersAdmin.POST("/:id/payment", s.addOrderPayment)
			ordersAdmin.POST("/:id/tracking", s.addOrderTracking)
			ordersAdmin.PUT("/:id/cancel", s.cancelOrder)
		}
	}

	// Supplier routes (admin/staff only)
	suppliers := v1.Group("/suppliers")
	suppliers.Use(s.authMiddleware(), s.staffMiddleware())
	{
		// Initialize supplier handler
		supplierHandler := NewSupplierHandler(s.supplierSvc, s.logger)
		
		// CRUD operations
		suppliers.GET("", supplierHandler.ListSuppliers)
		suppliers.POST("", supplierHandler.CreateSupplier)
		suppliers.GET("/:id", supplierHandler.GetSupplier)
		suppliers.PUT("/:id", supplierHandler.UpdateSupplier)
		suppliers.DELETE("/:id", supplierHandler.DeleteSupplier)
		
		// Adapter routes
		suppliers.GET("/adapters", supplierHandler.ListAdapters)
		suppliers.GET("/adapters/:name/capabilities", supplierHandler.GetAdapterCapabilities)
		suppliers.POST("/:id/test-connection", supplierHandler.TestAdapterConnection)
		
		// Sync routes
		suppliers.POST("/:id/sync/products", supplierHandler.SyncProducts)
		suppliers.POST("/:id/sync/inventory", supplierHandler.SyncInventory)
	}
	
	// POS (Point of Sale) routes (admin/staff only)
	pos := v1.Group("/pos")
	pos.Use(s.authMiddleware(), s.staffMiddleware())
	{
		// Order creation and processing
		pos.POST("/orders", s.createPOSOrder)
		pos.POST("/transactions/quick", s.processQuickPOSTransaction)
		
		// Inventory operations
		pos.POST("/inventory/check", s.checkPOSInventory)
		pos.POST("/inventory/reserve", s.reserveForPOSTransaction)
		pos.POST("/inventory/deduct", s.deductForDirectPOSSale)
		
		// In-store pickup
		pos.POST("/pickup/complete", s.completePickup)
	}
}

// Start starts the server
func (s *Server) Start() error {
	s.logger.Info("Starting REST server", zap.String("port", s.port))
	
	server := &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}
	
	return server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down REST server")
	
	server := &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}
	
	return server.Shutdown(ctx)
}

// loggerMiddleware creates a gin middleware for logging requests
func loggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		
		// Process request
		c.Next()
		
		// Log request details
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		
		logger.Info("HTTP Request",
			zap.String("path", path),
			zap.String("method", method),
			zap.Int("status", statusCode),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
		)
	}
}

// healthCheck returns a simple health check response
func (s *Server) healthCheck(c *gin.Context) {
	status := "ok"
	services := make(map[string]string)
	
	// Check service dependencies (basic connectivity)
	// Note: In production, you'd want to implement actual health checks for each service
	services["product"] = "ok"
	services["inventory"] = "ok"
	services["order"] = "ok"
	services["user"] = "ok"
	services["supplier"] = "ok"
	
	response := gin.H{
		"status":   status,
		"time":     time.Now().Format(time.RFC3339),
		"service":  "gateway",
		"version":  "1.0.0",
		"services": services,
	}
	
	c.JSON(http.StatusOK, response)
}

// respondWithError returns a formatted error response
func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"error": message,
	})
}

// respondWithSuccess returns a formatted success response
func respondWithSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{
		"success": true,
		"data":    data,
	})
}

// genericErrorHandler is a generic error handler
func genericErrorHandler(c *gin.Context, err error, logger *zap.Logger, operation string) {
	logger.Error("Operation failed",
		zap.String("operation", operation),
		zap.Error(err),
	)
	
	respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("%s failed: %v", operation, err))
}

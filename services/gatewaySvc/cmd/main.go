package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/gatewaySvc/internal/rest"
	"github.com/leonvanderhaeghen/stockplatform/services/gatewaySvc/internal/services"
)

// Config holds the application configuration
type Config struct {
	RestPort             string `env:"REST_PORT,default=8080"`
	JWTSecret            string `env:"JWT_SECRET,default=your-secret-key-here"`
	ProductServiceAddr   string `env:"PRODUCT_SERVICE_ADDR,default=localhost:50053"`
	InventoryServiceAddr string `env:"INVENTORY_SERVICE_ADDR,default=localhost:50054"`
	OrderServiceAddr     string `env:"ORDER_SERVICE_ADDR,default=localhost:50055"`
	UserServiceAddr      string `env:"USER_SERVICE_ADDR,default=localhost:50056"`
}

func main() {
	// Initialize logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting gateway service...")

	// Load configuration
	config := Config{
		RestPort:             "8080",
		JWTSecret:            "your-secret-key-here",
		ProductServiceAddr:   "localhost:50053",
		InventoryServiceAddr: "localhost:50054",
		OrderServiceAddr:     "localhost:50055",
		UserServiceAddr:      "localhost:50056",
	}
	
	// Check for environment variables
	if port := os.Getenv("REST_PORT"); port != "" {
		config.RestPort = port
	}
	
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.JWTSecret = jwtSecret
	}
	
	if addr := os.Getenv("PRODUCT_SERVICE_ADDR"); addr != "" {
		config.ProductServiceAddr = addr
	}
	
	if addr := os.Getenv("INVENTORY_SERVICE_ADDR"); addr != "" {
		config.InventoryServiceAddr = addr
	}
	
	if addr := os.Getenv("ORDER_SERVICE_ADDR"); addr != "" {
		config.OrderServiceAddr = addr
	}
	
	if addr := os.Getenv("USER_SERVICE_ADDR"); addr != "" {
		config.UserServiceAddr = addr
	}
	
	logger.Info("Configuration loaded", 
		zap.String("rest_port", config.RestPort),
		zap.String("product_service_addr", config.ProductServiceAddr),
		zap.String("inventory_service_addr", config.InventoryServiceAddr),
		zap.String("order_service_addr", config.OrderServiceAddr),
		zap.String("user_service_addr", config.UserServiceAddr),
	)

	// Initialize services
	productSvc, err := services.NewProductService(config.ProductServiceAddr, logger)
	if err != nil {
		logger.Fatal("Failed to create product service", zap.Error(err))
	}

	inventorySvc, err := services.NewInventoryService(config.InventoryServiceAddr, logger)
	if err != nil {
		logger.Fatal("Failed to create inventory service", zap.Error(err))
	}

	orderSvc, err := services.NewOrderService(config.OrderServiceAddr, logger)
	if err != nil {
		logger.Fatal("Failed to create order service", zap.Error(err))
	}

	userSvc, err := services.NewUserService(config.UserServiceAddr, logger)
	if err != nil {
		logger.Fatal("Failed to create user service", zap.Error(err))
	}

	// Initialize REST server
	server := rest.NewServer(
		productSvc,
		inventorySvc,
		orderSvc,
		userSvc,
		config.JWTSecret,
		config.RestPort,
		logger,
	)

	// Setup routes
	server.SetupRoutes()

	// Start server in a goroutine
	go func() {
		logger.Info("Starting REST server", zap.String("port", config.RestPort))
		if err := server.Start(); err != nil {
			logger.Fatal("Failed to start REST server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited properly")
}

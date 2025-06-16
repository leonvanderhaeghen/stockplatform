package main

import (
	"log"

	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/gatewaySvc/internal/config"
	"github.com/leonvanderhaeghen/stockplatform/services/gatewaySvc/internal/server"
)

func main() {
	// Initialize logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting gateway service...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	logger.Info("Configuration loaded",
		zap.String("port", cfg.Server.Port),
		zap.String("product_service", cfg.Services.ProductAddr),
		zap.String("inventory_service", cfg.Services.InventoryAddr),
		zap.String("order_service", cfg.Services.OrderAddr),
		zap.String("user_service", cfg.Services.UserAddr),
		zap.String("supplier_service", cfg.Services.SupplierAddr),
	)

	// Initialize server
	srv := server.New(cfg, logger)
	if err := srv.Initialize(); err != nil {
		logger.Fatal("Failed to initialize server", zap.Error(err))
	}

	// Start server (blocks until shutdown)
	if err := srv.Start(); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}

	logger.Info("Gateway service exited properly")
}

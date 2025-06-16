package main

import (
	"log"

	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/internal/config"
	"github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/internal/database"
	"github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/internal/server"
)

func main() {
	// Initialize logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting supplier service...")

	// Load configuration
	cfg := config.Load(logger)

	// Initialize database
	db, err := database.Initialize(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("Failed to close database", zap.Error(err))
		}
	}()

	// Initialize server
	srv := server.New(cfg, db, logger)
	if err := srv.Initialize(); err != nil {
		logger.Fatal("Failed to initialize server", zap.Error(err))
	}

	// Start server (blocks until shutdown)
	if err := srv.Start(); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}

	logger.Info("Supplier service exited properly")
}

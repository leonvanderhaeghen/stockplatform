package config

import (
	"os"

	"go.uber.org/zap"
)

// Config holds the application configuration
type Config struct {
	GRPCPort             string
	MongoURI             string
	Database             string
	ProductServiceAddr   string
	InventoryServiceAddr string
}

// Load loads configuration from environment variables
func Load(logger *zap.Logger) *Config {
	cfg := &Config{
		GRPCPort:             getEnv("GRPC_PORT", "50055"),
		MongoURI:             getEnv("MONGO_URI", "mongodb://localhost:27017"),
		Database:             getEnv("DATABASE_NAME", "stockplatform"),
		ProductServiceAddr:   getEnv("PRODUCT_SERVICE_ADDR", "product-service:50053"),
		InventoryServiceAddr: getEnv("INVENTORY_SERVICE_ADDR", "inventory-service:50054"),
	}

	logger.Info("Configuration loaded",
		zap.String("grpc_port", cfg.GRPCPort),
		zap.String("mongo_uri", maskSensitive(cfg.MongoURI)),
		zap.String("database", cfg.Database),
		zap.String("product_service_addr", cfg.ProductServiceAddr),
		zap.String("inventory_service_addr", cfg.InventoryServiceAddr),
	)

	return cfg
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// maskSensitive masks sensitive information for logging
func maskSensitive(value string) string {
	if len(value) > 20 {
		return value[:10] + "***" + value[len(value)-7:]
	}
	return "***"
}

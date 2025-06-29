package config

import (
	"os"

	"go.uber.org/zap"
)

// Config holds the application configuration
type Config struct {
	GRPCPort            string
	HTTPPort            string
	MongoURI            string
	Database            string
	SupplierServiceAddr string
	InventoryServiceAddr string
}

// Load loads configuration from environment variables with defaults
func Load(logger *zap.Logger) *Config {
	config := &Config{
		GRPCPort:            getEnvWithDefault("GRPC_PORT", "50053"),
		HTTPPort:            getEnvWithDefault("HTTP_PORT", "3001"),
		MongoURI:            getEnvWithDefault("MONGO_URI", "mongodb://localhost:27017"),
		Database:            getEnvWithDefault("DATABASE_NAME", "productdb"),
		SupplierServiceAddr: getEnvWithDefault("SUPPLIER_SERVICE_ADDR", "localhost:50057"),
		InventoryServiceAddr: getEnvWithDefault("INVENTORY_SERVICE_ADDR", "localhost:50052"),
	}

	// Log configuration (mask sensitive data)
	logger.Info("Configuration loaded",
		zap.String("grpc_port", config.GRPCPort),
		zap.String("http_port", config.HTTPPort),
		zap.String("mongo_uri", maskSensitiveData(config.MongoURI)),
		zap.String("database", config.Database),
		zap.String("supplier_service_addr", config.SupplierServiceAddr),
		zap.String("inventory_service_addr", config.InventoryServiceAddr),
	)

	return config
}

// getEnvWithDefault gets environment variable or returns default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// maskSensitiveData masks sensitive information in connection strings
func maskSensitiveData(data string) string {
	if len(data) > 20 {
		return data[:10] + "***" + data[len(data)-7:]
	}
	return "***"
}

package config

import (
	"os"

	"go.uber.org/zap"
)

// Config holds the application configuration
type Config struct {
	GRPCPort          string
	MongoURI          string
	Database          string
	OrderSvcURL       string
	DefaultLocationID string
}

// Load loads configuration from environment variables
func Load(logger *zap.Logger) *Config {
	cfg := &Config{
		GRPCPort:          getEnv("GRPC_PORT", "50054"),
		MongoURI:          getEnv("MONGO_URI", "mongodb://localhost:27017"),
		Database:          getEnv("DATABASE_NAME", "stockplatform"),
		OrderSvcURL:       getEnv("ORDER_SERVICE_URL", "order-service:50052"),
		DefaultLocationID: getEnv("DEFAULT_LOCATION_ID", "store-001"),
	}

	logger.Info("Configuration loaded",
		zap.String("grpc_port", cfg.GRPCPort),
		zap.String("mongo_uri", maskSensitive(cfg.MongoURI)),
		zap.String("database", cfg.Database),
		zap.String("order_service_url", cfg.OrderSvcURL),
		zap.String("default_location_id", cfg.DefaultLocationID),
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

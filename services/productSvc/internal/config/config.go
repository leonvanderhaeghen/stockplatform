package config

import (
	"os"

	"go.uber.org/zap"
)

// Config holds the application configuration
type Config struct {
	GRPCPort string
	HTTPPort string
	MongoURI string
	Database string
}

// Load loads configuration from environment variables with defaults
func Load(logger *zap.Logger) *Config {
	config := &Config{
		GRPCPort: getEnvWithDefault("GRPC_PORT", "50053"),
		HTTPPort: getEnvWithDefault("HTTP_PORT", "3001"),
		MongoURI: getEnvWithDefault("MONGO_URI", "mongodb://localhost:27017"),
		Database: getEnvWithDefault("DATABASE_NAME", "productdb"),
	}

	// Log configuration (mask sensitive data)
	logger.Info("Configuration loaded",
		zap.String("grpc_port", config.GRPCPort),
		zap.String("http_port", config.HTTPPort),
		zap.String("mongo_uri", maskSensitiveData(config.MongoURI)),
		zap.String("database", config.Database),
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

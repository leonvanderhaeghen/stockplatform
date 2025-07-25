package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the store service
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Services ServicesConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	URI      string
	Database string
}

// ServicesConfig holds configuration for other microservices
type ServicesConfig struct {
	ProductServiceAddr   string
	InventoryServiceAddr string
	OrderServiceAddr     string
	UserServiceAddr      string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("GRPC_PORT", getEnv("SERVER_PORT", "50058")),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			URI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
			Database: getEnv("DATABASE_NAME", "storedb"),
		},
		Services: ServicesConfig{
			ProductServiceAddr:   getEnv("PRODUCT_SERVICE_ADDR", "localhost:8081"),
			InventoryServiceAddr: getEnv("INVENTORY_SERVICE_ADDR", "localhost:8082"),
			OrderServiceAddr:     getEnv("ORDER_SERVICE_ADDR", "localhost:8083"),
			UserServiceAddr:      getEnv("USER_SERVICE_ADDR", "localhost:8084"),
		},
	}

	return cfg, nil
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

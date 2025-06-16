package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Services ServicesConfig `mapstructure:"services"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string `mapstructure:"port" validate:"required"`
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret string `mapstructure:"secret" validate:"required"`
}

// ServicesConfig holds service addresses
type ServicesConfig struct {
	ProductAddr   string `mapstructure:"product_addr" validate:"required"`
	InventoryAddr string `mapstructure:"inventory_addr" validate:"required"`
	OrderAddr     string `mapstructure:"order_addr" validate:"required"`
	UserAddr      string `mapstructure:"user_addr" validate:"required"`
	SupplierAddr  string `mapstructure:"supplier_addr" validate:"required"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level string `mapstructure:"level"`
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Set environment variable prefix
	viper.SetEnvPrefix("GATEWAY")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set defaults
	setDefaults()

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "8080")

	// JWT defaults
	viper.SetDefault("jwt.secret", "your-secret-key-here")

	// Service defaults
	viper.SetDefault("services.product_addr", "localhost:50053")
	viper.SetDefault("services.inventory_addr", "localhost:50054")
	viper.SetDefault("services.order_addr", "localhost:50055")
	viper.SetDefault("services.user_addr", "localhost:50056")
	viper.SetDefault("services.supplier_addr", "localhost:50057")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
}

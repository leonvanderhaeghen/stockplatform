package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds logger configuration
type Config struct {
	Level       string `mapstructure:"level"`
	Environment string `mapstructure:"environment"`
	ServiceName string `mapstructure:"service_name"`
}

// New creates a new structured logger based on configuration
func New(cfg Config) (*zap.Logger, error) {
	// Set default values
	if cfg.Level == "" {
		cfg.Level = "info"
	}
	if cfg.Environment == "" {
		cfg.Environment = "development"
	}
	if cfg.ServiceName == "" {
		cfg.ServiceName = "unknown"
	}

	// Parse log level
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	// Configure encoder
	var encoderConfig zapcore.EncoderConfig
	if strings.ToLower(cfg.Environment) == "production" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Configure core
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	// Create logger with service name field
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	logger = logger.With(zap.String("service", cfg.ServiceName))

	return logger, nil
}

// NewDevelopment creates a development logger
func NewDevelopment(serviceName string) (*zap.Logger, error) {
	return New(Config{
		Level:       "debug",
		Environment: "development",
		ServiceName: serviceName,
	})
}

// NewProduction creates a production logger
func NewProduction(serviceName string) (*zap.Logger, error) {
	return New(Config{
		Level:       "info",
		Environment: "production",
		ServiceName: serviceName,
	})
}

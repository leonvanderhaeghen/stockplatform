package application

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new configured logger
func NewLogger(env string) (*zap.Logger, error) {
	var config zap.Config
	
	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	
	return config.Build()
}

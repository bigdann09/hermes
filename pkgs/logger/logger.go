package logger

import (
	"github.com/bigdann09/notifications/internal/config"
	"go.uber.org/zap"
)

func NewLogger(cfg *config.AppConfig) *zap.Logger {
	var level zap.AtomicLevel
	if cfg.Environment == "development" {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	config := zap.Config{
		Level:            level,
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewProductionEncoderConfig(),
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()
	return logger
}

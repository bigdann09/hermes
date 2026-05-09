package services

import (
	"github.com/bigdann09/notifications/internal/config"
	"github.com/bigdann09/notifications/internal/database"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	Logger   *zap.Logger
	Config   *config.Config
	Database *gorm.DB
}

func NewService(logger *zap.Logger, cfg *config.Config) *Service {
	logger.Info("registering database service...")
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	logger.Info("database connected")

	return &Service{
		Logger:   logger,
		Config:   cfg,
		Database: db,
	}
}

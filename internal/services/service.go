package services

import (
	"github.com/bigdann09/notifications/internal/config"
	"github.com/bigdann09/notifications/internal/infrastructure/cache"
	"github.com/bigdann09/notifications/internal/infrastructure/database"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	Logger   *zap.Logger
	Config   *config.Config
	Database *gorm.DB
	Cache    *cache.Cache
}

func NewService(logger *zap.Logger, cfg *config.Config) *Service {
	logger.Info("registering database service...")
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	logger.Info("database connected")

	logger.Info("registering cache service...")
	cache := cache.NewCache(&cfg.Cache)
	logger.Info("cache connected")

	return &Service{
		Logger:   logger,
		Config:   cfg,
		Database: db,
		Cache:    cache,
	}
}

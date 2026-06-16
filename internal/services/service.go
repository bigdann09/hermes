package services

import (
	"github.com/bigdann09/notifications/internal/config"
	"github.com/bigdann09/notifications/internal/infrastructure/cache"
	"github.com/bigdann09/notifications/internal/infrastructure/database"
	"github.com/bigdann09/notifications/internal/infrastructure/kafka"
	"github.com/bigdann09/notifications/pkgs/validators"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	Logger        *zap.Logger
	Config        *config.Config
	Database      *gorm.DB
	Cache         *cache.Cache
	KafkaProducer *kafka.KafkaProducer
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

	logger.Info("registering valitators...")
	validator := validators.NewValidator()
	validator.Register(
		&validators.ValidationField{
			Tag:         "has_notification_type",
			Func:        validator.HasNotificationType,
			Translation: "{0} is not a supported notification type",
		},
	)
	logger.Info("validators registered")

	logger.Info("setting up kafka producer...")
	kafka_producer, err := kafka.NewKafkaProducer(&cfg.Kafka)
	if err != nil {
		logger.Fatal("failed to connect to kafka", zap.Error(err))
	}
	logger.Info("kafka connected")

	return &Service{
		Logger:        logger,
		Config:        cfg,
		Database:      db,
		Cache:         cache,
		KafkaProducer: kafka_producer,
	}
}

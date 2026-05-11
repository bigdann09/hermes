package seeders

import (
	"fmt"

	"github.com/bigdann09/notifications/internal/config"
	"github.com/bigdann09/notifications/internal/infrastructure/database"
	"github.com/bigdann09/notifications/pkgs/logger"
	"go.uber.org/zap"
)

func Seed() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	logger := logger.NewLogger(&cfg.App)
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return
	}

	SeedUsers(logger, db)
	SeedNotificationPreferences(logger, db)
	logger.Info("Seed completed successfully")
}

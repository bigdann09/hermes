package seeders

import (
	"github.com/bigdann09/notifications/internal/models"
	"github.com/bigdann09/notifications/internal/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SeedNotificationPreferences(logger *zap.Logger, db *gorm.DB) {
	userRepository := repositories.NewUserRepository(db)
	repository := repositories.NewNotificationPreferenceRepository(db)

	users, err := userRepository.FindAll()
	if err != nil {
		logger.Error("failed to load users for notification preferences", zap.Error(err))
		return
	}

	for _, user := range users {
		preference := models.NotificationPreference{
			UserID:        user.ID,
			EmailEnabled:  true,
			PushEnabled:   true,
			SocketEnabled: true,
		}

		existing, _ := repository.FindByUserID(user.ID)
		if existing.ID == "" {
			logger.Info("Creating notification preference", zap.String("user_id", user.ID))
			repository.Create(&preference)
			continue
		}

		logger.Info("Updating notification preference", zap.String("user_id", user.ID))
		preference.ID = existing.ID
		repository.Update(&preference)
	}
}

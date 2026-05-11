package seeders

import (
	"github.com/bigdann09/notifications/internal/models"
	"github.com/bigdann09/notifications/internal/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var NOTIFICATION_PREFERENCE_DATA = []models.NotificationPreference{
	{
		UserID:        "59608cd4-893a-489c-98f8-941f781e2187",
		EmailEnabled:  true,
		PushEnabled:   true,
		SocketEnabled: true,
	},
}

func SeedNotificationPreferences(logger *zap.Logger, db *gorm.DB) {
	repository := repositories.NewNotificationPreferenceRepository(db)
	for _, notification_preference := range NOTIFICATION_PREFERENCE_DATA {
		existing_notification_preference, _ := repository.FindByUserID(notification_preference.UserID)
		if existing_notification_preference.ID == "" {
			logger.Info("Creating notification preference", zap.String("user_id", notification_preference.UserID))
			repository.Create(&notification_preference)
		} else {
			logger.Info("Updating notification preference", zap.String("user_id", notification_preference.UserID))
			notification_preference.ID = existing_notification_preference.ID
			repository.Update(&notification_preference)
		}
	}
}

package seeders

import (
	"github.com/bigdann09/notifications/internal/models"
	"github.com/bigdann09/notifications/internal/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var USER_DATA = []models.User{
	{
		Username: "John Doe",
		Email:    "johndoe@gmail.com",
	},
}

func SeedUsers(logger *zap.Logger, db *gorm.DB) {
	repository := repositories.NewUserRepository(db)
	for _, user := range USER_DATA {
		existing_user, _ := repository.FindByEmail(user.Email)
		if existing_user.ID == "" {
			logger.Info("Creating user", zap.String("email", user.Email))
			repository.Create(&user)
		} else {
			logger.Info("Updating user", zap.String("email", user.Email))
			user.ID = existing_user.ID
			repository.Update(&user)
		}
	}
}

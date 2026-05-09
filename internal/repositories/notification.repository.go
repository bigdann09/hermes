package repositories

import (
	"github.com/bigdann09/notifications/internal/models"
	"gorm.io/gorm"
)

type INotificationRepository interface {
}

type NotificationRepository struct {
	db    *gorm.DB
	table string
}

func NewNotificationRepository(db *gorm.DB) INotificationRepository {
	return &NotificationRepository{db: db, table: "notifications"}
}

func (repository *NotificationRepository) FindAll() (*[]models.Notification, error) {
	var notifications *[]models.Notification
	result := repository.db.Table(repository.table).Find(&notifications)
	if result.Error != nil {
		return nil, result.Error
	}
	return notifications, nil
}

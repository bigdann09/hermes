package repositories

import "gorm.io/gorm"

type INotificationRepository interface {
}

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) INotificationRepository {
	return &NotificationRepository{db}
}

func (repository *NotificationRepository) FindAll() {

}

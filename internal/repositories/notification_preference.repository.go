package repositories

import (
	"github.com/bigdann09/notifications/internal/models"
	"gorm.io/gorm"
)

type NotificationPreferenceRepository interface {
	Create(notification_preference *models.NotificationPreference) error
	Update(notification_preference *models.NotificationPreference) error
	Delete(notification_preference *models.NotificationPreference) error
	FindByID(id string) (models.NotificationPreference, error)
	FindByUserID(user_id string) (models.NotificationPreference, error)
	FindAll() ([]models.NotificationPreference, error)
}

type notificationPreferenceRepository struct {
	db *gorm.DB
}

func NewNotificationPreferenceRepository(db *gorm.DB) NotificationPreferenceRepository {
	return &notificationPreferenceRepository{db: db}
}

func (r *notificationPreferenceRepository) Create(notification_preference *models.NotificationPreference) error {
	return r.db.Create(notification_preference).Error
}

func (r *notificationPreferenceRepository) Update(notification_preference *models.NotificationPreference) error {
	return r.db.Save(notification_preference).Error
}

func (r *notificationPreferenceRepository) Delete(notification_preference *models.NotificationPreference) error {
	return r.db.Delete(notification_preference).Error
}

func (r *notificationPreferenceRepository) FindByID(id string) (models.NotificationPreference, error) {
	var notification_preference models.NotificationPreference
	err := r.db.Where("id = ?", id).First(&notification_preference).Error
	return notification_preference, err
}

func (r *notificationPreferenceRepository) FindByUserID(user_id string) (models.NotificationPreference, error) {
	var notification_preference models.NotificationPreference
	err := r.db.Where("user_id = ?", user_id).First(&notification_preference).Error
	return notification_preference, err
}

func (r *notificationPreferenceRepository) FindAll() ([]models.NotificationPreference, error) {
	var notification_preferences []models.NotificationPreference
	err := r.db.Find(&notification_preferences).Error
	return notification_preferences, err
}

package notification

import (
	"github.com/bigdann09/notifications/internal/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type INotificationService interface {
	Send() error
}

type NotificationService struct {
	logger     *zap.Logger
	repository repositories.INotificationRepository
}

func NewNotificationService(database *gorm.DB, logger *zap.Logger) INotificationService {
	return &NotificationService{
		repository: repositories.NewNotificationRepository(database),
		logger:     logger,
	}
}

func (s *NotificationService) Send() error {
	return nil
}

package notification

import (
	"time"

	"github.com/bigdann09/notifications/internal/dtos"
	"github.com/bigdann09/notifications/internal/infrastructure/cache"
	"github.com/bigdann09/notifications/internal/models"
	"github.com/bigdann09/notifications/internal/repositories"
	"github.com/bigdann09/notifications/pkgs/apiresponse"
	"github.com/bigdann09/notifications/pkgs/pagination"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type INotificationService interface {
	Send() error
	FindAll(query dtos.NotificationQuery) (*pagination.Pagination[models.Notification], *apiresponse.APIResponse)
}

type NotificationService struct {
	logger     *zap.Logger
	repository repositories.INotificationRepository
	cache      *cache.Cache
}

func NewNotificationService(database *gorm.DB, logger *zap.Logger, cache *cache.Cache) INotificationService {
	return &NotificationService{
		repository: repositories.NewNotificationRepository(database),
		logger:     logger,
		cache:      cache,
	}
}

func (srv *NotificationService) FindAll(query dtos.NotificationQuery) (*pagination.Pagination[models.Notification], *apiresponse.APIResponse) {
	cache_key := "notifications:all"
	var cached pagination.Pagination[models.Notification]
	if err := srv.cache.Get(cache_key, &cached); err == nil {
		srv.logger.Info("notifications retrieved successfully from cache")
		return &cached, nil
	}

	result, err := srv.repository.FindAllPaginated(&query)
	if err != nil {
		srv.logger.Error("could not retrieve notifications", zap.Error(err))
		return nil, apiresponse.InternalServerError("could not retrieve notifications")
	}

	srv.logger.Info("notifications retrieved successfully")
	srv.cache.Set(cache_key, result, 10*time.Minute)
	return result, nil
}

func (s *NotificationService) Send() error {
	return nil
}

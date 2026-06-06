package notification

import (
	"time"

	"github.com/bigdann09/notifications/internal/dtos"
	"github.com/bigdann09/notifications/internal/infrastructure/cache"
	"github.com/bigdann09/notifications/internal/infrastructure/kafka"
	"github.com/bigdann09/notifications/internal/models"
	"github.com/bigdann09/notifications/internal/repositories"
	"github.com/bigdann09/notifications/internal/services/notification/channels"
	"github.com/bigdann09/notifications/pkgs/apiresponse"
	"github.com/bigdann09/notifications/pkgs/pagination"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const notificationTopic = "notifications"

type INotificationService interface {
	Send(request dtos.NotificationRequest) *apiresponse.APIResponse
	FindAll(query dtos.NotificationQuery) (*pagination.Pagination[models.Notification], *apiresponse.APIResponse)
}

type NotificationService struct {
	logger                  *zap.Logger
	notification_repository repositories.INotificationRepository
	user_repository         repositories.IUserRepository
	preference_repository   repositories.NotificationPreferenceRepository
	kafka_producer          *kafka.KafkaProducer
	cache                   *cache.Cache
}

func NewNotificationService(
	database *gorm.DB,
	logger *zap.Logger,
	cache *cache.Cache,
	kafkaProducer *kafka.KafkaProducer,
) INotificationService {
	return &NotificationService{
		notification_repository: repositories.NewNotificationRepository(database),
		user_repository:         repositories.NewUserRepository(database),
		preference_repository:   repositories.NewNotificationPreferenceRepository(database),
		logger:                  logger,
		cache:                   cache,
		kafka_producer:          kafkaProducer,
	}
}

func (srv *NotificationService) FindAll(query dtos.NotificationQuery) (*pagination.Pagination[models.Notification], *apiresponse.APIResponse) {
	cacheKey := "notifications:all"
	var cached pagination.Pagination[models.Notification]
	if err := srv.cache.Get(cacheKey, &cached); err == nil {
		srv.logger.Info("notifications retrieved successfully from cache")
		return &cached, nil
	}

	result, err := srv.notification_repository.FindAllPaginated(&query)
	if err != nil {
		srv.logger.Error("could not retrieve notifications", zap.Error(err))
		return nil, apiresponse.InternalServerError("could not retrieve notifications")
	}

	srv.logger.Info("notifications retrieved successfully")
	srv.cache.Set(cacheKey, result, 10*time.Minute)
	return result, nil
}

func (srv *NotificationService) Send(request dtos.NotificationRequest) *apiresponse.APIResponse {
	user, err := srv.user_repository.FindByID(request.UserID)
	if err != nil {
		srv.logger.Error("user not found", zap.Error(err))
		return apiresponse.NotFound("user not found")
	}

	preferences, err := srv.preference_repository.FindByUserID(request.UserID)
	if err != nil {
		srv.logger.Warn("notification preferences not found, using defaults", zap.Error(err))
		preferences = models.NotificationPreference{
			EmailEnabled: true,
		}
	}

	channel_list := make([]string, 0, 3)
	channel_list = append(channel_list, string(channels.Database))
	if preferences.EmailEnabled {
		channel_list = append(channel_list, string(channels.Email))
	}
	if preferences.PushEnabled {
		channel_list = append(channel_list, string(channels.Push))
	}
	if preferences.SocketEnabled {
		channel_list = append(channel_list, string(channels.Websocket))
	}

	message := dtos.KafkaNotificationMessage{
		UserID:   request.UserID,
		Email:    user.Email,
		Title:    request.Title,
		Message:  request.Message,
		Type:     string(request.Type),
		Metadata: request.Metadata,
		Channels: channel_list,
	}

	if err := srv.kafka_producer.Publish(notificationTopic, request.UserID, message); err != nil {
		srv.logger.Error("failed to publish notification", zap.Error(err))
		return apiresponse.InternalServerError("failed to queue notification")
	}

	srv.cache.Delete("notifications:all")
	srv.logger.Info("notification queued successfully", zap.String("user_id", request.UserID))
	return apiresponse.NewWithData(202, "notification queued successfully", message)
}

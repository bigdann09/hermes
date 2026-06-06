package channels

import (
	"encoding/json"

	"github.com/bigdann09/notifications/internal/models"
	"github.com/bigdann09/notifications/internal/repositories"
)

type DatabaseChannel struct {
	repository repositories.INotificationRepository
}

func NewDatabaseChannel(repository repositories.INotificationRepository) *DatabaseChannel {
	return &DatabaseChannel{repository: repository}
}

func (channel *DatabaseChannel) Type() string {
	return string(Database)
}

func (channel *DatabaseChannel) Send(payload SendNotificationPayload) error {
	metadata := payload.Data
	if metadata == nil {
		metadata = map[string]any{}
	}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	notification := &models.Notification{
		UserID:   payload.UserID,
		Type:     models.NotificationType(payload.Type),
		Title:    payload.Title,
		Metadata: string(metadataJSON),
	}
	return channel.repository.Create(notification)
}

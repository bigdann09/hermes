package notification

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/bigdann09/notifications/internal/dtos"
	"github.com/bigdann09/notifications/internal/services/notification/channels"
	"go.uber.org/zap"
)

type DispatcherService struct {
	channels []channels.IChannel
	logger   *zap.Logger
}

func NewDispatcherService(logger *zap.Logger, channel_list []channels.IChannel) *DispatcherService {
	return &DispatcherService{
		channels: channel_list,
		logger:   logger,
	}
}

func (dispatcher *DispatcherService) Dispatch(msg *sarama.ConsumerMessage) {
	var payload dtos.KafkaNotificationMessage
	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		dispatcher.logger.Error("failed to unmarshal notification message", zap.Error(err))
		return
	}

	dispatcher.logger.Info(
		"dispatching notification",
		zap.String("user_id", payload.UserID),
		zap.Strings("channels", payload.Channels),
	)

	send_payload := channels.SendNotificationPayload{
		UserID:  payload.UserID,
		Email:   payload.Email,
		Title:   payload.Title,
		Message: payload.Message,
		Type:    payload.Type,
		Data:    payload.Metadata,
	}

	for _, channelName := range payload.Channels {
		send_payload.Channels = []channels.Channel{channels.Channel(channelName)}
		dispatcher.Send(send_payload)
	}
}

func (dispatcher *DispatcherService) Send(payload channels.SendNotificationPayload) {
	for _, channel := range payload.Channels {
		for _, registered := range dispatcher.channels {
			if registered.Type() != string(channel) {
				continue
			}
			if err := registered.Send(payload); err != nil {
				dispatcher.logger.Error(
					"failed to send notification",
					zap.String("channel", string(channel)),
					zap.Error(err),
				)
			}
		}
	}
}

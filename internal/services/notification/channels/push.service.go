package channels

import (
	"go.uber.org/zap"
)

type PushChannel struct {
	logger *zap.Logger
}

func NewPushChannel(logger *zap.Logger) *PushChannel {
	return &PushChannel{logger: logger}
}

func (channel *PushChannel) Type() string {
	return string(Push)
}

func (channel *PushChannel) Send(payload SendNotificationPayload) error {
	channel.logger.Info(
		"push notification sent",
		zap.String("user_id", payload.UserID),
		zap.String("title", payload.Title),
	)
	return nil
}

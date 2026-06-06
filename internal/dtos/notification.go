package dtos

import "github.com/bigdann09/notifications/internal/models"

type NotificationRequest struct {
	UserID   string                  `json:"user_id" binding:"required,uuid"`
	Type     models.NotificationType `json:"type" binding:"required,has_notification_type"`
	Title    string                  `json:"title" binding:"required"`
	Message  string                  `json:"message" binding:"required"`
	Metadata map[string]any          `json:"metadata,omitempty"`
}

type KafkaNotificationMessage struct {
	UserID   string         `json:"user_id"`
	Email    string         `json:"email"`
	Title    string         `json:"title"`
	Message  string         `json:"message"`
	Type     string         `json:"type"`
	Metadata map[string]any `json:"metadata,omitempty"`
	Channels []string       `json:"channels"`
}

type NotificationQuery struct {
	Page   uint                    `form:"page,omitempty"`
	Limit  uint                    `form:"limit,omitempty"`
	UserID string                  `form:"user_id" binding:"required,uuid"`
	Type   models.NotificationType `form:"type,omitempty" binding:"has_notification_type"`
	IsRead bool                    `form:"is_read,omitempty" binding:"boolean"`
	SortBy string                  `form:"order_by,omitempty"`
	Order  string                  `form:"order,omitempty"`
}

func (query *NotificationQuery) GetOffset() uint {
	return (query.Page - 1) * query.Limit
}

func (query *NotificationQuery) Default() {
	if query.Page == 0 {
		query.Page = 1
	}

	if query.Limit == 0 {
		query.Limit = 10
	}

	if query.SortBy == "" {
		query.SortBy = "created_at"
	}

	if query.Order == "" {
		query.Order = "desc"
	}
}

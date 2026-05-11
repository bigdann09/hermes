package dtos

import "github.com/bigdann09/notifications/internal/models"

type NotificationQuery struct {
	Page   uint                    `form:"page,omitempty"`
	Limit  uint                    `form:"limit,omitempty"`
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

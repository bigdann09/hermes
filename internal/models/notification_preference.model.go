package models

import "time"

type NotificationPreference struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID        string    `gorm:"uniqueIndex;not null" json:"user_id"`
	EmailEnabled  bool      `gorm:"default:true" json:"email_enabled"`
	PushEnabled   bool      `gorm:"default:true" json:"push_enabled"`
	SocketEnabled bool      `gorm:"default:true" json:"socket_enabled"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

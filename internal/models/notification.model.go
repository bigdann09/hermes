package models

import "time"

type NotificationType string

var (
	System NotificationType = "system"
	Social NotificationType = "social"
)

type Notification struct {
	ID        string           `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    string           `gorm:"index;not null" json:"user_id"`
	Type      NotificationType `gorm:"not null" json:"type"`
	Title     string           `gorm:"not null" json:"title"`
	Metadata  string           `gorm:"type:text" json:"metadata"`
	ReadAt    time.Time        `gorm:"null" json:"read_at"`
	CreatedAt time.Time        `gorm:"autoCreateTime" json:"created_at"`
}

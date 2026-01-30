package models

import "time"

type Notification struct {
	ID         int
	CategoryID int
	ItemName   string
	ItemURL    string
	SentAt     *time.Time
}

type PendingNotification struct {
	ID         int
	CategoryID int
	ItemName   string
	ItemURL    string
	UserID     int // для MarkDelivered(notification_id, user_id)
	TelegramID int64
}

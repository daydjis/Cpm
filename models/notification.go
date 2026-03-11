package models

import "time"

type Notification struct {
	ID         int
	CategoryID int
	ItemName   string
	ItemURL    string // ссылка на страницу товара
	ImageURL   string // ссылка на картинку товара
	SentAt     *time.Time
}

type PendingNotification struct {
	ID           int
	CategoryID   int
	CategoryName string
	ItemName     string
	ItemURL      string // ссылка на страницу товара
	ImageURL     string // ссылка на картинку товара
	UserID       int    // для MarkDelivered(notification_id, user_id)
	TelegramID   int64
}

package repository

import (
	"awesomeProject3/models"
	"database/sql"
	"log"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) GetPending() ([]models.PendingNotification, error) {
	rows, err := r.db.Query(`
        SELECT n.id, n.category_id, n.item_name, n.item_url, s.user_id, u.telegram_id
        FROM notifications n
        JOIN subscriptions s ON s.category_id = n.category_id AND s.filter_signature = n.filter_signature
        JOIN users u ON u.id = s.user_id
        LEFT JOIN notification_deliveries d ON d.notification_id = n.id AND d.user_id = s.user_id
        WHERE d.id IS NULL
        ORDER BY n.id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.PendingNotification
	for rows.Next() {
		var n models.PendingNotification
		if err := rows.Scan(&n.ID, &n.CategoryID, &n.ItemName, &n.ItemURL, &n.UserID, &n.TelegramID); err != nil {
			log.Println("Scan error:", err)
			continue
		}
		list = append(list, n)
	}
	return list, rows.Err()
}

func (r *NotificationRepository) MarkDelivered(notificationID, userID int) error {
	_, err := r.db.Exec(`
		INSERT INTO notification_deliveries (notification_id, user_id) VALUES ($1, $2)
		ON CONFLICT (notification_id, user_id) DO NOTHING
	`, notificationID, userID)
	return err
}

func (r *NotificationRepository) CreateOrIgnore(categoryID int, itemName, itemURL, filterSignature string) error {
	_, err := r.db.Exec(`
		INSERT INTO notifications (category_id, item_name, item_url, filter_signature, sent_at)
		VALUES ($1, $2, $3, $4, NULL)
		ON CONFLICT (category_id, item_url, filter_signature) DO NOTHING
	`, categoryID, itemName, itemURL, filterSignature)
	return err
}

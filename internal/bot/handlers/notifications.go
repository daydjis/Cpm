package handlers

import (
	"awesomeProject3/models"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendNewNotifications() {
	if svc == nil {
		log.Println("services not set")
		return
	}
	notifications, err := svc.Notification.GetPending()
	if err != nil {
		log.Println("Failed to fetch notifications:", err)
		return
	}
	for _, n := range notifications {
		if err := sendNotification(n); err != nil {
			log.Println("Failed to send notification:", err)
		}
	}
}

func sendNotification(n models.PendingNotification) error {
	log.Printf("Sending to %d: %s\n", n.TelegramID, n.ItemName)
	msg := tgbotapi.NewMessage(n.TelegramID, n.ItemName+"\n"+n.ItemURL)
	if botAPI != nil {
		if _, err := botAPI.Send(msg); err != nil {
			return fmt.Errorf("telegram send error: %w", err)
		}
	}
	if err := svc.Notification.MarkDelivered(n.ID, n.UserID); err != nil {
		return fmt.Errorf("update sent_at error: %w", err)
	}
	log.Println("Marked as sent:", n.ID)
	return nil
}

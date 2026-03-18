package handlers

import (
	"awesomeProject3/models"
	"fmt"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const maxItemsPerMessage = 20

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

	userNotifs := make(map[int64][]models.PendingNotification)
	for _, n := range notifications {
		userNotifs[n.TelegramID] = append(userNotifs[n.TelegramID], n)
	}

	for chatID, list := range userNotifs {
		if len(list) == 0 {
			continue
		}

		catMap := make(map[string][]models.PendingNotification)
		for _, n := range list {
			name := n.CategoryName
			if name == "" {
				name = "Без категории"
			}
			catMap[name] = append(catMap[name], n)
		}

		for catName, catList := range catMap {
			if len(catList) == 0 {
				continue
			}

			if len(catList) == 1 && catList[0].ImageURL != "" {
				if err := sendSingleWithPhoto(chatID, catList[0]); err != nil {
					log.Println("Failed to send notification photo:", err)
					continue
				}
			} else {
				for i := 0; i < len(catList); i += maxItemsPerMessage {
					end := i + maxItemsPerMessage
					if end > len(catList) {
						end = len(catList)
					}
					if err := sendCategoryListMessage(chatID, catName, catList[i:end]); err != nil {
						log.Println("Failed to send notification batch:", err)
						break
					}
				}
			}

			for _, n := range catList {
				if err := svc.Notification.MarkDelivered(n.ID, n.UserID); err != nil {
					log.Println("Failed to mark delivered:", err)
				}
			}
		}
	}
}

func sendSingleWithPhoto(chatID int64, n models.PendingNotification) error {
	if botAPI == nil {
		return fmt.Errorf("bot API not initialized")
	}

	title, price := splitNameAndPrice(n.ItemName)

	caption := fmt.Sprintf(
		"<b><a href=\"%s\">%s</a></b>\nЦена: <b>%s</b>",
		htmlEscape(n.ItemURL),
		htmlEscape(title),
		htmlEscape(price),
	)

	resp, err := http.Get(n.ImageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	photo := tgbotapi.FileReader{
		Name:   "image.webp",
		Reader: resp.Body,
	}

	msg := tgbotapi.NewPhoto(chatID, photo)
	msg.Caption = caption
	msg.ParseMode = tgbotapi.ModeHTML

	log.Printf("Sending photo to %d: %s\n", chatID, n.ItemName)

	_, err = botAPI.Send(msg)
	return err
}

func sendCategoryListMessage(chatID int64, categoryName string, list []models.PendingNotification) error {
	if botAPI == nil {
		return fmt.Errorf("bot API not initialized")
	}

	var b strings.Builder
	b.WriteString("🔔 <b>Новые предложения по подписке</b>\n")
	b.WriteString("📂 <b>" + htmlEscape(categoryName) + "</b>\n\n")

	for _, n := range list {
		title, price := splitNameAndPrice(n.ItemName)
		line := fmt.Sprintf(
			"• <a href=\"%s\">%s</a> — <b>%s</b>\n",
			htmlEscape(n.ItemURL),
			htmlEscape(title),
			htmlEscape(price),
		)
		b.WriteString(line)
	}

	msg := tgbotapi.NewMessage(chatID, b.String())
	msg.ParseMode = tgbotapi.ModeHTML
	// Чтобы не спамить превьюшками на каждую ссылку.
	msg.DisableWebPagePreview = true

	log.Printf("Sending batch to %d: %d items\n", chatID, len(list))

	if _, err := botAPI.Send(msg); err != nil {
		return fmt.Errorf("telegram send batch error: %w", err)
	}
	return nil
}

func splitNameAndPrice(s string) (title, price string) {
	parts := strings.SplitN(s, " | ", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return s, ""
}

func htmlEscape(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
	)
	return replacer.Replace(s)
}

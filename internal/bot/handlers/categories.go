package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleCategories(message *tgbotapi.Message) {
	if svc == nil {
		log.Println("services not set")
		return
	}
	showCategoryKeyboard(message.Chat.ID, "Доступные категории (нажмите для подписки):", "subscribe_")
}

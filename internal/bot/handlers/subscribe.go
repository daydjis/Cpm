package handlers

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleSubscribe(message *tgbotapi.Message) {
	if svc == nil {
		log.Println("services not set")
		return
	}
	showCategoryKeyboard(message.Chat.ID, "Выберите категорию для подписки:", "subscribe_")
}

func showCategoryKeyboard(chatID int64, text, callbackPrefix string) {
	categories, err := svc.Category.List()
	if err != nil {
		log.Println("List categories:", err)
		msg := tgbotapi.NewMessage(chatID, "Не удалось загрузить категории.")
		botSend(msg)
		return
	}
	if len(categories) == 0 {
		msg := tgbotapi.NewMessage(chatID, "Пока нет доступных категорий. Добавьте через /add_category.")
		botSend(msg)
		return
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(categories))
	for _, c := range categories {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(c.Name+" ("+c.Slug+")", callbackPrefix+strconv.Itoa(c.ID)),
		))
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	botSend(msg)
}

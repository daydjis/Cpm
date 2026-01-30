package handlers

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleUnsubscribe(message *tgbotapi.Message) {
	if svc == nil {
		log.Println("services not set")
		return
	}
	subs, err := svc.Subscription.GetUserSubscriptions(message.Chat.ID)
	if err != nil {
		log.Println("GetUserSubscriptions:", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "Ошибка загрузки подписок.")
		botSend(msg)
		return
	}
	if len(subs) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "У вас пока нет подписок. Используйте «Подписаться» или /subscribe.")
		botSend(msg)
		return
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(subs))
	for _, s := range subs {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ "+s.CategoryName, "unsubscribe_"+strconv.Itoa(s.CategoryID)),
		))
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите категорию, от которой отписаться:")
	msg.ReplyMarkup = keyboard
	botSend(msg)
}

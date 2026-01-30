package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleStart(message *tgbotapi.Message) {
	if svc == nil {
		log.Println("services not set")
		return
	}
	if err := svc.User.RegisterOrIgnore(message.Chat.ID, message.Chat.UserName); err != nil {
		log.Println("Error inserting user:", err)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Привет! Вы зарегистрированы для рассылки.\n\nИспользуйте кнопки меню или команды ниже.")
	msg.ReplyMarkup = MainMenuKeyboard()
	botSend(msg)
}

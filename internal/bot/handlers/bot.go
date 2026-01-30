package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var botAPI *tgbotapi.BotAPI

func SetBot(b *tgbotapi.BotAPI) {
	botAPI = b
}

func botSend(c tgbotapi.Chattable) {
	if botAPI != nil {
		_, err := botAPI.Send(c)
		if err != nil {
			return
		}

	}
}

func botRequest(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	if botAPI == nil {
		return nil, nil
	}
	return botAPI.Request(c)
}

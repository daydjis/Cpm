package bot

import (
	"awesomeProject3/internal/config"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

func Init() {
	if config.TelegramToken == "" {
		log.Fatal("TG_BOT_TOKEN or TELEGRAM_TOKEN is not set")
	}
	var err error
	Bot, err = tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		log.Panic(err)
	}
	Bot.Debug = false
	log.Printf("Authorized on account %s", Bot.Self.UserName)
}

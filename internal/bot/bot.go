package bot

import (
	"awesomeProject3/internal/config"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/proxy"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

func SetBot(b *tgbotapi.BotAPI) {
	Bot = b
}

func Init() {
	if config.TelegramToken == "" {
		log.Fatal("TG_BOT_TOKEN or TELEGRAM_TOKEN is not set")
	}

	// SOCKS5 proxy
	dialer, err := proxy.SOCKS5(
		"tcp",
		"127.0.0.1:10808",
		nil,
		proxy.Direct,
	)

	if err != nil {
		log.Fatalf("proxy error: %v", err)
	}

	transport := &http.Transport{}
	transport.Dial = dialer.Dial

	client := &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second,
	}

	for i := 0; i < 5; i++ {
		Bot, err = tgbotapi.NewBotAPIWithClient(
			config.TelegramToken,
			tgbotapi.APIEndpoint,
			client,
		)

		if err == nil {
			break
		}

		log.Printf("Telegram init failed (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to init bot after retries: %v", err)
	}

	Bot.Debug = false
	log.Printf("Authorized on account %s", Bot.Self.UserName)
}

func Stop() {
	if Bot != nil {
		Bot.StopReceivingUpdates()
	}
}

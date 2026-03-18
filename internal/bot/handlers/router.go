package handlers

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUpdates() {
	if botAPI == nil {
		log.Println("handlers: bot not set, call SetBot first")
		return
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			handleCallback(update.CallbackQuery)
			continue
		}
		if update.Message == nil {
			continue
		}

		msg := update.Message
		chatID := msg.Chat.ID
		text := msg.Text
		st := getState(chatID)
		if strings.HasPrefix(st, "filter_") {
			parts := strings.SplitN(st, "_", 4) // filter_<subID>_<step>
			if len(parts) >= 3 {
				if subID, err := strconv.Atoi(parts[1]); err == nil {
					if handleFilterStep(chatID, subID, parts[2], text) {
						continue
					}
				}
			}
		}
		if st == stateAddCategorySlug {
			handleAddCategoryText(chatID, text)
			continue
		}

		cmd := msg.Command()

		switch {
		case cmd == "start" || text == "/start":
			handleStart(msg)
		case cmd == "subscribe" || text == "📂 Подписаться":
			handleSubscribe(msg)
		case cmd == "unsubscribe" || text == "📋 Мои подписки":
			handleUnsubscribe(msg)
		case cmd == "categories" || text == "📁 Все категории":
			handleCategories(msg)
		case cmd == "add_category" || text == "➕ Добавить категорию":
			handleAddCategory(msg)
		default:
			out := tgbotapi.NewMessage(chatID, "Неизвестная команда.\n\n"+HelpText)
			botSend(out)
		}
	}
}

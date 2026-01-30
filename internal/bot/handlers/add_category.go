package handlers

import (
	"awesomeProject3/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const stateAddCategorySlug = "add_category_slug"

func handleAddCategory(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите slug категории (как в URL сайта).\nПримеры: eldar, orks, space-marines\n\nСсылка будет: .../warhammer-40000/<slug>")
	botSend(msg)
	setState(message.Chat.ID, stateAddCategorySlug)
}

func handleAddCategoryText(chatID int64, text string) {
	slug := trimLower(text)
	if slug == "" {
		msg := tgbotapi.NewMessage(chatID, "Slug не может быть пустым.")
		botSend(msg)
		clearState(chatID)
		return
	}
	name := slug
	if err := svc.Category.Create(name, slug); err != nil {
		if err == service.ErrSlugExists {
			msg := tgbotapi.NewMessage(chatID, "Категория с таким slug уже есть.")
			botSend(msg)
		} else {
			msg := tgbotapi.NewMessage(chatID, "Ошибка: "+err.Error())
			botSend(msg)
		}
		clearState(chatID)
		return
	}
	clearState(chatID)
	msg := tgbotapi.NewMessage(chatID, "✅ Категория «"+name+"» добавлена. Теперь можно подписаться через «Подписаться».")
	botSend(msg)
}

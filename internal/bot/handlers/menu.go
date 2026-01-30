package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Todo доделать логику Ui
func MainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📂 Подписаться"),
			tgbotapi.NewKeyboardButton("📋 Мои подписки"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📁 Все категории"),
			tgbotapi.NewKeyboardButton("➕ Добавить категорию"),
		),
	)
}

const (
	HelpText = `Команды:
/start — регистрация и меню
/subscribe — подписаться на категорию (после подписки можно настроить фильтры: текст, цена, регион, типы размещения)
/unsubscribe — отписаться от категории
/categories — список категорий
/add_category — добавить категорию (slug: eldar, orks и т.д.)

Фильтры формируют URL для парсинга: поиск, цена от/до, регион, типы (0,1=все; 3,5=мастерские и литейные).`
)

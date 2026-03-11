package handlers

import (
	"awesomeProject3/internal/filter"
	"awesomeProject3/internal/service"
	"awesomeProject3/models"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleCallback(callback *tgbotapi.CallbackQuery) {
	if svc == nil {
		return
	}
	data := callback.Data
	switch {
	case strings.HasPrefix(data, "subscribe_"):
		handleSubscribeCallback(callback, strings.TrimPrefix(data, "subscribe_"))
	case strings.HasPrefix(data, "unsubscribe_"):
		handleUnsubscribeCallback(callback, strings.TrimPrefix(data, "unsubscribe_"))
	case strings.HasPrefix(data, "filters_"):
		handleFiltersChoiceCallback(callback, data)
	case strings.HasPrefix(data, "filter_skip_"):
		handleFilterSkipCallback(callback, data)
	case strings.HasPrefix(data, "filter_pro_"):
		handleFilterProCallback(callback, data)
	default:
		answerCallback(callback.ID, "Неизвестное действие")
	}
}

func handleSubscribeCallback(callback *tgbotapi.CallbackQuery, idStr string) {
	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		answerCallback(callback.ID, "Ошибка")
		return
	}
	telegramID := callback.From.ID
	subID, err := svc.Subscription.SubscribeByCategoryID(telegramID, categoryID)
	if err != nil {
		if err == service.ErrCategoryNotFound {
			answerCallback(callback.ID, "Категория не найдена")
		} else if err == service.ErrUserNotFound {
			answerCallback(callback.ID, "Сначала отправьте /start")
		} else {
			log.Println("Subscribe error:", err)
			answerCallback(callback.ID, "Ошибка подписки")
		}
		return
	}
	answerCallback(callback.ID, "Вы подписались!")
	edit := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "✅ Подписка оформлена. Настроить фильтры поиска?")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да", "filters_yes_"+strconv.Itoa(subID)),
			tgbotapi.NewInlineKeyboardButtonData("Позже", "filters_no_"+strconv.Itoa(subID)),
		),
	)

	edit.ReplyMarkup = &keyboard
	botSend(edit)
}

func handleFiltersChoiceCallback(callback *tgbotapi.CallbackQuery, data string) {
	// filters_yes_<subID> или filters_no_<subID>
	parts := strings.SplitN(data, "_", 3)
	if len(parts) != 3 {
		answerCallback(callback.ID, "Ошибка")
		return
	}
	subID, err := strconv.Atoi(parts[2])
	if err != nil {
		answerCallback(callback.ID, "Ошибка")
		return
	}
	chatID := callback.Message.Chat.ID
	answerCallback(callback.ID, "")
	edit := tgbotapi.NewEditMessageText(chatID, callback.Message.MessageID, "✅ Подписка оформлена.")
	edit.ReplyMarkup = nil
	botSend(edit)
	if parts[1] == "yes" {
		askFilterSearch(chatID, subID)
	}
}

func handleFilterSkipCallback(callback *tgbotapi.CallbackQuery, suffix string) {
	// filter_skip_<subID>_<step>
	parts := strings.SplitN(suffix, "_", 4)

	if len(parts) != 4 {
		answerCallback(callback.ID, "Ошибка")
		return
	}
	subID, err := strconv.Atoi(parts[2])
	if err != nil {
		answerCallback(callback.ID, "Ошибка")
		return
	}
	chatID := callback.Message.Chat.ID
	answerCallback(callback.ID, "Пропущено")
	clearState(chatID)
	handleFilterSkip(chatID, subID, parts[3])
}

func handleFilterProCallback(callback *tgbotapi.CallbackQuery, suffix string) {
	// filter_pro_<subID>_<types>
	parts := strings.SplitN(suffix, "_", 4)
	if len(parts) != 4 {
		answerCallback(callback.ID, "Ошибка")
		return
	}
	subID, err := strconv.Atoi(parts[2])
	if err != nil {
		answerCallback(callback.ID, "Ошибка")
		return
	}
	typesStr := parts[3]
	types := filter.ParseProTypes(typesStr)
	setFilterParam(subID, func(p *models.FilterParams) { p.ProTypes = types })
	chatID := callback.Message.Chat.ID
	answerCallback(callback.ID, "Выбрано")
	clearState(chatID)
	finishFilterFlow(chatID, subID)
}

func handleUnsubscribeCallback(callback *tgbotapi.CallbackQuery, idStr string) {
	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		answerCallback(callback.ID, "Ошибка")
		return
	}
	telegramID := callback.From.ID
	if err := svc.Subscription.Unsubscribe(telegramID, categoryID); err != nil {
		if err == service.ErrUserNotFound {
			answerCallback(callback.ID, "Сначала отправьте /start")
		} else {
			log.Println("Unsubscribe error:", err)
			answerCallback(callback.ID, "Ошибка отписки")
		}
		return
	}
	answerCallback(callback.ID, "Вы отписались.")
	edit := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "✅ Вы отписались от выбранной категории.")
	edit.ReplyMarkup = nil
	botSend(edit)
}

func answerCallback(callbackID, text string) {
	answer := tgbotapi.NewCallback(callbackID, text)
	botRequest(answer)
}

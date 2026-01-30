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
	case len(data) >= 10 && data[:10] == "subscribe_":
		handleSubscribeCallback(callback, data[10:])
	case len(data) >= 12 && data[:12] == "unsubscribe_":
		handleUnsubscribeCallback(callback, data[12:])
	case len(data) >= 14 && data[:11] == "filters_yes":
		handleFiltersChoiceCallback(callback, data[:14])

	case len(data) >= 13 && data[:12] == "filter_skip_":
		handleFilterSkipCallback(callback, data[11:])
	case len(data) >= 10 && data[:10] == "filter_pro":
		handleFilterProCallback(callback, data[10:])
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

func handleFiltersChoiceCallback(callback *tgbotapi.CallbackQuery, suffix string) {

	parts := strings.SplitN(suffix, "_", 3)
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
	parts := strings.SplitN(suffix, "_", 3)

	if len(parts) != 3 {
		answerCallback(callback.ID, "Ошибка")
		return
	}
	subID, err := strconv.Atoi(parts[1])
	if err != nil {
		answerCallback(callback.ID, "Ошибка")
		return
	}
	chatID := callback.Message.Chat.ID
	answerCallback(callback.ID, "Пропущено")
	clearState(chatID)
	handleFilterSkip(chatID, subID, parts[2])
}

func handleFilterProCallback(callback *tgbotapi.CallbackQuery, suffix string) {
	idx := strings.Index(suffix, "_")
	if idx < 0 {
		answerCallback(callback.ID, "Ошибка1")
		return
	}
	subID, err := strconv.Atoi(strconv.Itoa(idx))
	if err != nil {
		answerCallback(callback.ID, "Ошибка2")
		return
	}
	typesStr := suffix[idx+1:]
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

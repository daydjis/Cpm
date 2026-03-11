package handlers

import (
	"awesomeProject3/internal/filter"
	"awesomeProject3/models"
	"log"
	"strconv"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	stepFilterSearch   = "search"
	stepFilterPriceMin = "price_min"
	stepFilterPriceMax = "price_max"
	stepFilterRegion   = "region"
	stepFilterProTypes = "pro_types"
	stepFilterDone     = "done"
)

var filterParamsMu sync.Mutex
var filterParams = make(map[int]models.FilterParams)

func askFilterSearch(chatID int64, subID int) {
	setState(chatID, "filter_"+strconv.Itoa(subID)+"_"+stepFilterSearch)
	msg := tgbotapi.NewMessage(chatID, "🔍 Введите текст поиска (или нажмите Пропустить):")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Пропустить", "filter_skip_"+strconv.Itoa(subID)+"_search")),
	)
	botSend(msg)
}

func askFilterPriceMin(chatID int64, subID int) {
	setState(chatID, "filter_"+strconv.Itoa(subID)+"_"+stepFilterPriceMin)
	msg := tgbotapi.NewMessage(chatID, "💰 Цена от (руб, число). 0 = без ограничения:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Пропустить", "filter_skip_"+strconv.Itoa(subID)+"_price_min")),
	)
	botSend(msg)
}

func askFilterPriceMax(chatID int64, subID int) {
	setState(chatID, "filter_"+strconv.Itoa(subID)+"_"+stepFilterPriceMax)
	msg := tgbotapi.NewMessage(chatID, "💰 Цена до (руб, число). 0 = без ограничения:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Пропустить", "filter_skip_"+strconv.Itoa(subID)+"_price_max")),
	)
	botSend(msg)
}

func askFilterRegion(chatID int64, subID int) {
	setState(chatID, "filter_"+strconv.Itoa(subID)+"_"+stepFilterRegion)
	msg := tgbotapi.NewMessage(chatID, "📍 Регион (id с сайта, 0 = все). Или Пропустить:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Пропустить", "filter_skip_"+strconv.Itoa(subID)+"_region")),
	)
	botSend(msg)
}

func askFilterProTypes(chatID int64, subID int) {
	setState(chatID, "filter_"+strconv.Itoa(subID)+"_"+stepFilterProTypes)
	msg := tgbotapi.NewMessage(chatID, "🏷 Типы размещения на сайте:\n• 0 — Пользователи \n• 1 - Магазины \n• 3 - Мастерские \n• 5 -  литейные\n• Или свой набор через запятую: 0,1,3,5")
	rows := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Пользователи", "filter_pro_"+strconv.Itoa(subID)+"_0"),
			tgbotapi.NewInlineKeyboardButtonData("Магазины", "filter_pro_"+strconv.Itoa(subID)+"_1"),
			tgbotapi.NewInlineKeyboardButtonData("Мастерские", "filter_pro_"+strconv.Itoa(subID)+"_3"),
			tgbotapi.NewInlineKeyboardButtonData("Литейные", "filter_pro_"+strconv.Itoa(subID)+"_5"),
		),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Пропустить (без фильтра)", "filter_skip_"+strconv.Itoa(subID)+"_pro_types")),
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	botSend(msg)
}

func getFilterParams(subID int) models.FilterParams {
	filterParamsMu.Lock()
	defer filterParamsMu.Unlock()
	p, ok := filterParams[subID]
	if !ok {
		return models.FilterParams{}
	}
	return p
}

func setFilterParam(subID int, set func(*models.FilterParams)) {
	filterParamsMu.Lock()
	defer filterParamsMu.Unlock()
	p := filterParams[subID]
	set(&p)
	filterParams[subID] = p
}

func clearFilterParams(subID int) {
	filterParamsMu.Lock()
	defer filterParamsMu.Unlock()
	delete(filterParams, subID)
}

func finishFilterFlow(chatID int64, subID int) {
	p := getFilterParams(subID)
	if err := svc.Subscription.UpdateFilters(subID, p); err != nil {
		log.Println("UpdateFilters:", err)
		botSend(tgbotapi.NewMessage(chatID, "Ошибка сохранения фильтров."))
	} else {
		botSend(tgbotapi.NewMessage(chatID, "✅ Фильтры сохранены. Новые объявления по этим параметрам будут приходить в рассылке."))
	}
	clearFilterParams(subID)
	clearState(chatID)
}

func handleFilterStep(chatID int64, subID int, step, text string) bool {
	switch step {
	case stepFilterSearch:
		setFilterParam(subID, func(p *models.FilterParams) { p.SearchText = strings.TrimSpace(text) })
		askFilterPriceMin(chatID, subID)
		return true
	case stepFilterPriceMin:
		v, _ := strconv.Atoi(strings.TrimSpace(text))
		setFilterParam(subID, func(p *models.FilterParams) { p.PriceMin = v })
		askFilterPriceMax(chatID, subID)
		return true
	case stepFilterPriceMax:
		v, _ := strconv.Atoi(strings.TrimSpace(text))
		setFilterParam(subID, func(p *models.FilterParams) { p.PriceMax = v })
		askFilterRegion(chatID, subID)
		return true
	case stepFilterRegion:
		v, _ := strconv.Atoi(strings.TrimSpace(text))
		setFilterParam(subID, func(p *models.FilterParams) { p.RegionID = v })
		askFilterProTypes(chatID, subID)
		return true
	case stepFilterProTypes:
		types := filter.ParseProTypes(text)
		setFilterParam(subID, func(p *models.FilterParams) { p.ProTypes = types })
		clearState(chatID)
		finishFilterFlow(chatID, subID)
		return true
	}
	return false
}

func handleFilterSkip(chatID int64, subID int, step string) {
	switch step {
	case "search":
		askFilterPriceMin(chatID, subID)
	case "price_min":
		askFilterPriceMax(chatID, subID)
	case "price_max":
		askFilterRegion(chatID, subID)
	case "region":
		askFilterProTypes(chatID, subID)
	case "pro_types":
		clearState(chatID)
		finishFilterFlow(chatID, subID)
	}
}

package bot

import (
	handlers "awesomeProject3/internal/bot/handlers"
	"awesomeProject3/internal/service"
)

func SetServices(s *service.Services) {
	handlers.SetServices(s)
}

func SetHandlersBot() {
	handlers.SetBot(Bot)
}

func HandleUpdates() {
	handlers.HandleUpdates()
}

func SendNewNotifications() {
	handlers.SendNewNotifications()
}

package handlers

import "awesomeProject3/internal/service"

var svc *service.Services

func SetServices(s *service.Services) {
	svc = s
}

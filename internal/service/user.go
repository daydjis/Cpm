package service

import (
	"awesomeProject3/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) RegisterOrIgnore(telegramID int64, username string) error {
	return s.userRepo.CreateOrIgnore(telegramID, username)
}

func (s *UserService) GetUserIDByTelegramID(telegramID int64) (int, error) {
	return s.userRepo.GetIDByTelegramID(telegramID)
}

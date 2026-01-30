package service

import (
	"awesomeProject3/internal/repository"
	"awesomeProject3/models"
)

type NotificationService struct {
	notifRepo *repository.NotificationRepository
}

func NewNotificationService(notifRepo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{notifRepo: notifRepo}
}

func (s *NotificationService) GetPending() ([]models.PendingNotification, error) {
	return s.notifRepo.GetPending()
}

func (s *NotificationService) MarkDelivered(notificationID, userID int) error {
	return s.notifRepo.MarkDelivered(notificationID, userID)
}

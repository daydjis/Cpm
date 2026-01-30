package service

import "awesomeProject3/internal/repository"

type Services struct {
	User         *UserService
	Category     *CategoryService
	Subscription *SubscriptionService
	Notification *NotificationService
}

func NewServicesFromRepos(
	userRepo *repository.UserRepository,
	categoryRepo *repository.CategoryRepository,
	subRepo *repository.SubscriptionRepository,
	notifRepo *repository.NotificationRepository,
) *Services {
	return &Services{
		User:         NewUserService(userRepo),
		Category:     NewCategoryService(categoryRepo),
		Subscription: NewSubscriptionService(subRepo, userRepo, categoryRepo),
		Notification: NewNotificationService(notifRepo),
	}
}

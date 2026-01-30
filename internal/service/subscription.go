package service

import (
	"awesomeProject3/internal/filter"
	"awesomeProject3/internal/repository"
	"awesomeProject3/models"
	"database/sql"
	"errors"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrUserNotFound     = errors.New("user not found")
)

type SubscriptionService struct {
	subRepo      *repository.SubscriptionRepository
	userRepo     *repository.UserRepository
	categoryRepo *repository.CategoryRepository
}

func NewSubscriptionService(
	subRepo *repository.SubscriptionRepository,
	userRepo *repository.UserRepository,
	categoryRepo *repository.CategoryRepository,
) *SubscriptionService {
	return &SubscriptionService{
		subRepo:      subRepo,
		userRepo:     userRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *SubscriptionService) SubscribeByTelegramAndCategoryName(telegramID int64, categoryName string) (subID int, err error) {
	category, err := s.categoryRepo.GetByName(categoryName)
	if err != nil {
		return 0, ErrCategoryNotFound
	}
	return s.SubscribeByCategoryID(telegramID, category.ID)
}

func (s *SubscriptionService) SubscribeByCategoryID(telegramID int64, categoryID int) (subID int, err error) {
	if _, err := s.categoryRepo.GetByID(categoryID); err != nil {
		return 0, ErrCategoryNotFound
	}
	userID, err := s.userRepo.GetIDByTelegramID(telegramID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrUserNotFound
		}
		return 0, err
	}
	return s.subRepo.Subscribe(userID, categoryID)
}

func (s *SubscriptionService) UpdateFilters(subID int, p models.FilterParams) error {
	sig := filter.BuildFilterSignature(p)
	proTypes := filter.FormatProTypes(p.ProTypes)
	return s.subRepo.UpdateFilters(subID, p.SearchText, p.PriceMin, p.PriceMax, p.RegionID, proTypes, sig)
}

func (s *SubscriptionService) GetSubscriptionByID(subID int) (*models.Subscription, error) {
	return s.subRepo.GetByID(subID)
}

func (s *SubscriptionService) ListForScheduler() ([]repository.SubscriptionWithCategory, error) {
	return s.subRepo.ListForScheduler()
}

func (s *SubscriptionService) GetUserSubscriptions(telegramID int64) ([]repository.UserSubscriptionView, error) {
	userID, err := s.userRepo.GetIDByTelegramID(telegramID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s.subRepo.ListByUser(userID)
}

func (s *SubscriptionService) Unsubscribe(telegramID int64, categoryID int) error {
	userID, err := s.userRepo.GetIDByTelegramID(telegramID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}
		return err
	}
	return s.subRepo.Unsubscribe(userID, categoryID)
}

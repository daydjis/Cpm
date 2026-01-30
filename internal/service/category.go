package service

import (
	"awesomeProject3/internal/repository"
	"awesomeProject3/models"
	"errors"
	"regexp"
	"strings"
)

var ErrSlugExists = errors.New("категория с таким slug уже есть")

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (s *CategoryService) List() ([]models.Category, error) {
	return s.categoryRepo.List()
}

// SlugRegex — допустимые символы в slug: латиница, цифры, дефис.
var SlugRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9\-]*[a-z0-9]$|^[a-z0-9]$`)

func (s *CategoryService) Create(name, slug string) error {
	slug = strings.TrimSpace(strings.ToLower(slug))
	if slug == "" {
		return errors.New("slug не может быть пустым")
	}
	if !SlugRegex.MatchString(slug) {
		return errors.New("slug: только латиница, цифры и дефис (например eldar, space-marines)")
	}
	inserted, err := s.categoryRepo.Create(name, slug)
	if err != nil {
		return err
	}
	if !inserted {
		return ErrSlugExists
	}
	return nil
}

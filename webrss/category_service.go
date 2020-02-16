package webrss

import (
	"fmt"
	"time"

	"github.com/Alkemic/webrss/repository"
)

type categoryRepository interface {
	List(...string) ([]repository.Category, error)
	Get(int64) (repository.Category, error)
	Delete(int64) error
	Update(repository.Category) error
	GetNextByOrder(repository.Category) (repository.Category, error)
	GetPrevByOrder(repository.Category) (repository.Category, error)
	Create(repository.Category) error
	SelectMaxOrder() (int, error)
}

type feedRepository interface {
	Get(id int64) (repository.Feed, error)
	ListForCategories([]int64) ([]repository.Feed, error)
	Update(entry repository.Feed) error
}

type CategoryService struct {
	categoryRepository categoryRepository
	feedRepository     feedRepository
}

func NewCategoryService(categoryRepository categoryRepository, feedRepository feedRepository) *CategoryService {
	return &CategoryService{
		categoryRepository: categoryRepository,
		feedRepository:     feedRepository,
	}
}

func (s CategoryService) Get(id int64) (repository.Category, error) {
	entries, err := s.categoryRepository.Get(id)
	if err != nil {
		return repository.Category{}, fmt.Errorf("error fetching category: %w", err)
	}
	return entries, nil
}

func (s CategoryService) List(params ...string) ([]repository.Category, error) {
	categories, err := s.categoryRepository.List(params...)
	if err != nil {
		return nil, fmt.Errorf("error fetching categories: %w", err)
	}
	ids := make([]int64, 0, len(categories))
	for i, category := range categories {
		ids = append(ids, category.ID)
		categories[i].Feeds = make([]repository.Feed, 0)
	}

	feeds, err := s.feedRepository.ListForCategories(ids)
	if err != nil {
		return nil, fmt.Errorf("error fetching feeds for categories: %w", err)
	}
	for _, feed := range feeds {
		for i, category := range categories {
			if category.ID == feed.CategoryID {
				categories[i].Feeds = append(categories[i].Feeds, feed)
			}
		}
	}
	return categories, nil
}

func (s CategoryService) Create(category repository.Category) error {
	category.CreatedAt.Time = time.Now()
	maxOrder, err := s.categoryRepository.SelectMaxOrder()
	if err != nil {
		return fmt.Errorf("error fetching max order for category: %w", err)
	}
	category.Order = maxOrder + 1
	if err := s.categoryRepository.Create(category); err != nil {
		return fmt.Errorf("error creating category: %w", err)
	}
	return nil
}

func (s CategoryService) Delete(id int64) error {
	err := s.categoryRepository.Delete(id)
	if err != nil {
		return fmt.Errorf("error fetching category: %w", err)
	}
	return nil
}

func (s CategoryService) Update(category repository.Category) error {
	err := s.categoryRepository.Update(category)
	if err != nil {
		return fmt.Errorf("error fetching category: %w", err)
	}
	return nil
}

func (s CategoryService) MoveUp(id int64) error {
	entry, err := s.categoryRepository.Get(id)
	if err != nil {
		return fmt.Errorf("error fetching category: %w", err)
	}
	nextEntry, err := s.categoryRepository.GetNextByOrder(entry)
	if err != nil {
		return fmt.Errorf("error fetching next category: %w", err)
	}

	nextEntry.Order, entry.Order = entry.Order, nextEntry.Order
	if err := s.categoryRepository.Update(entry); err != nil {
		return fmt.Errorf("error updating category: %w", err)
	}
	if err := s.categoryRepository.Update(nextEntry); err != nil {
		return fmt.Errorf("error updating next category: %w", err)
	}

	return nil
}

func (s CategoryService) MoveDown(id int64) error {
	entry, err := s.categoryRepository.Get(id)
	if err != nil {
		return fmt.Errorf("error fetching category: %w", err)
	}
	prevEntry, err := s.categoryRepository.GetPrevByOrder(entry)
	if err != nil {
		return fmt.Errorf("error fetching previous category: %w", err)
	}

	prevEntry.Order, entry.Order = entry.Order, prevEntry.Order
	if err := s.categoryRepository.Update(entry); err != nil {
		return fmt.Errorf("error updating category: %w", err)
	}
	if err := s.categoryRepository.Update(prevEntry); err != nil {
		return fmt.Errorf("error updating previous category: %w", err)
	}

	return nil
}

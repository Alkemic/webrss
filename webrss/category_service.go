package webrss

import (
	"context"
	"fmt"
	"time"

	"github.com/Alkemic/webrss/repository"
)

type categoryRepository interface {
	List(ctx context.Context, params ...string) ([]repository.Category, error)
	Get(ctx context.Context, id int64) (repository.Category, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, category repository.Category) error
	GetNextByOrder(ctx context.Context, category repository.Category) (repository.Category, error)
	GetPrevByOrder(ctx context.Context, category repository.Category) (repository.Category, error)
	Create(ctx context.Context, category repository.Category) error
	SelectMaxOrder(ctx context.Context) (int, error)
}

type feedRepository interface {
	Get(ctx context.Context, id int64) (repository.Feed, error)
	Create(ctx context.Context, feed repository.Feed) (int64, error)
	ListForCategories(ctx context.Context, ids []int64) ([]repository.Feed, error)
	List(ctx context.Context) ([]repository.Feed, error)
	Update(ctx context.Context, entry repository.Feed) error
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
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

func (s CategoryService) Get(ctx context.Context, id int64) (repository.Category, error) {
	entries, err := s.categoryRepository.Get(ctx, id)
	if err != nil {
		return repository.Category{}, fmt.Errorf("error fetching category: %w", err)
	}
	return entries, nil
}

func (s CategoryService) List(ctx context.Context, params ...string) ([]repository.Category, error) {
	categories, err := s.categoryRepository.List(ctx, params...)
	if err != nil {
		return nil, fmt.Errorf("error fetching categories: %w", err)
	}
	ids := make([]int64, 0, len(categories))
	for i, category := range categories {
		ids = append(ids, category.ID)
		categories[i].Feeds = make([]repository.Feed, 0)
	}

	feeds, err := s.feedRepository.ListForCategories(ctx, ids)
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

func (s CategoryService) Create(ctx context.Context, category repository.Category) error {
	category.CreatedAt.Time = time.Now()
	maxOrder, err := s.categoryRepository.SelectMaxOrder(ctx)
	if err != nil {
		return fmt.Errorf("error fetching max order for category: %w", err)
	}
	category.Order = maxOrder + 1
	if err := s.categoryRepository.Create(ctx, category); err != nil {
		return fmt.Errorf("error creating category: %w", err)
	}
	return nil
}

func (s CategoryService) Delete(ctx context.Context, id int64) error {
	err := s.categoryRepository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error fetching category: %w", err)
	}
	return nil
}

func (s CategoryService) Update(ctx context.Context, category repository.Category) error {
	err := s.categoryRepository.Update(ctx, category)
	if err != nil {
		return fmt.Errorf("error fetching category: %w", err)
	}
	return nil
}

func (s CategoryService) MoveUp(ctx context.Context, id int64) error {
	entry, err := s.categoryRepository.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("error fetching category: %w", err)
	}
	nextEntry, err := s.categoryRepository.GetNextByOrder(ctx, entry)
	if err != nil {
		return fmt.Errorf("error fetching next category: %w", err)
	}

	nextEntry.Order, entry.Order = entry.Order, nextEntry.Order
	if err := s.categoryRepository.Update(ctx, entry); err != nil {
		return fmt.Errorf("error updating category: %w", err)
	}
	if err := s.categoryRepository.Update(ctx, nextEntry); err != nil {
		return fmt.Errorf("error updating next category: %w", err)
	}

	return nil
}

func (s CategoryService) MoveDown(ctx context.Context, id int64) error {
	entry, err := s.categoryRepository.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("error fetching category: %w", err)
	}
	prevEntry, err := s.categoryRepository.GetPrevByOrder(ctx, entry)
	if err != nil {
		return fmt.Errorf("error fetching previous category: %w", err)
	}

	prevEntry.Order, entry.Order = entry.Order, prevEntry.Order
	if err := s.categoryRepository.Update(ctx, entry); err != nil {
		return fmt.Errorf("error updating category: %w", err)
	}
	if err := s.categoryRepository.Update(ctx, prevEntry); err != nil {
		return fmt.Errorf("error updating previous category: %w", err)
	}

	return nil
}

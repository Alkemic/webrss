package webrss

import (
	"context"
	"fmt"
	"time"

	"github.com/Alkemic/webrss/repository"
)

func (s WebRSSService) GetCategory(ctx context.Context, id int64) (repository.Category, error) {
	entries, err := s.categoryRepository.Get(ctx, id)
	if err != nil {
		return repository.Category{}, fmt.Errorf("error fetching category: %w", err)
	}
	return entries, nil
}

func (s WebRSSService) ListCategories(ctx context.Context, params ...string) ([]repository.Category, error) {
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

func (s WebRSSService) CreateCategory(ctx context.Context, category repository.Category) error {
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

func (s WebRSSService) DeleteCategory(ctx context.Context, id int64) error {
	category, err := s.GetCategory(ctx, id)
	if err != nil {
		return fmt.Errorf("cannot fetch catory for delete: %w", err)
	}
	category.DeletedAt = repository.NewNullTime(s.nowFn())
	if err := s.UpdateCategory(ctx, category); err != nil {
		return fmt.Errorf("error deleting category: %w", err)
	}
	return nil
}

func (s WebRSSService) UpdateCategory(ctx context.Context, category repository.Category) error {
	category.UpdatedAt = repository.NewNullTime(s.nowFn())
	err := s.categoryRepository.Update(ctx, category)
	if err != nil {
		return fmt.Errorf("error fetching category: %w", err)
	}
	return nil
}

func (s WebRSSService) MoveCategoryDown(ctx context.Context, id int64) error {
	entry, err := s.categoryRepository.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("error fetching category: %w", err)
	}
	nextEntry, err := s.categoryRepository.GetNextByOrder(ctx, entry.Order)
	if err != nil {
		return fmt.Errorf("error fetching next category: %w", err)
	}

	nextEntry.Order, entry.Order = entry.Order, nextEntry.Order
	if err := s.transactionRepository.Begin(ctx); err != nil {
		return fmt.Errorf("cannot start transation when moving down category: %w", err)
	}
	defer s.transactionRepository.Rollback(ctx)
	if err := s.UpdateCategory(ctx, entry); err != nil {
		return fmt.Errorf("error updating category: %w", err)
	}
	if err := s.UpdateCategory(ctx, nextEntry); err != nil {
		return fmt.Errorf("error updating next category: %w", err)
	}
	if err := s.transactionRepository.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit transation when moving down category: %w", err)
	}

	return nil
}

func (s WebRSSService) MoveCategoryUp(ctx context.Context, id int64) error {
	entry, err := s.categoryRepository.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("error fetching category: %w", err)
	}
	prevEntry, err := s.categoryRepository.GetPrevByOrder(ctx, entry.Order)
	if err != nil {
		return fmt.Errorf("error fetching previous category: %w", err)
	}

	prevEntry.Order, entry.Order = entry.Order, prevEntry.Order
	if err := s.transactionRepository.Begin(ctx); err != nil {
		return fmt.Errorf("cannot start transation when moving up category: %w", err)
	}
	defer s.transactionRepository.Rollback(ctx)
	if err := s.UpdateCategory(ctx, entry); err != nil {
		return fmt.Errorf("error updating category: %w", err)
	}
	if err := s.UpdateCategory(ctx, prevEntry); err != nil {
		return fmt.Errorf("error updating previous category: %w", err)
	}
	if err := s.transactionRepository.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit transation when moving up category: %w", err)
	}

	return nil
}

package webrss

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Alkemic/webrss/repository"

	"github.com/Alkemic/webrss/feed_fetcher"
)

type feedFetcher interface {
	Fetch(ctx context.Context, url string) (feed_fetcher.Feed, error)
}

type FeedService struct {
	nowFn                 func() time.Time
	feedRepository        feedRepository
	entryRepository       entryRepository
	transactionRepository transactionRepository
	feedFetcher           feedFetcher
}

func NewFeedService(feedRepository feedRepository, entryRepository entryRepository, transactionRepository transactionRepository, feedFetcher feedFetcher) *FeedService {
	return &FeedService{
		nowFn:                 time.Now,
		feedRepository:        feedRepository,
		entryRepository:       entryRepository,
		transactionRepository: transactionRepository,
		feedFetcher:           feedFetcher,
	}
}

func (s FeedService) Get(ctx context.Context, id int64) (repository.Feed, error) {
	feed, err := s.feedRepository.Get(ctx, id)
	if err != nil {
		return repository.Feed{}, fmt.Errorf("error getting feed: %w", err)
	}
	return feed, nil
}

func (s FeedService) Create(ctx context.Context, feedURL string, categoryID int64) error {
	feeder, err := s.feedFetcher.Fetch(ctx, feedURL)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}
	feed := feeder.Feed(ctx)
	entries := feeder.Entries(ctx)
	now := repository.NewTime(s.nowFn())
	feed.CategoryID = categoryID
	feed.CreatedAt = now
	feed.LastReadAt = repository.NewTime(time.Date(1900, 1, 1, 1, 1, 1, 1, time.UTC))

	if err := s.transactionRepository.Begin(ctx); err != nil {
		return fmt.Errorf("cannot start transation when creating new feed: %w", err)
	}
	defer s.transactionRepository.Rollback(ctx)

	feedID, err := s.feedRepository.Create(ctx, feed)
	if err != nil {
		return fmt.Errorf("error creating new feed: %w", err)
	}
	if err := s.SaveEntries(ctx, feedID, entries); err != nil {
		return fmt.Errorf("error saving entries: %w", err)
	}
	if err := s.transactionRepository.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit transation when creating new feed: %w", err)
	}
	return nil
}

func updateEntry(a, b repository.Entry) repository.Entry {
	a.Author = b.Author
	a.Summary = b.Summary
	a.Title = b.Title
	a.PublishedAt = b.PublishedAt
	return a
}

func (s FeedService) SaveEntries(ctx context.Context, feedID int64, entries []repository.Entry) error {
	now := repository.NewTime(s.nowFn())
	for _, entry := range entries {
		existingEntry, err := s.entryRepository.GetByURL(ctx, entry.Link)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("error fetching entry: %w", err)
		} else if errors.Is(err, sql.ErrNoRows) {
			entry.FeedID = feedID
			entry.CreatedAt = now
			if err := s.entryRepository.Create(ctx, entry); err != nil {
				return fmt.Errorf("error creating entry: %w", err)
			}
		} else {
			entry = updateEntry(existingEntry, entry)
			entry.UpdatedAt = repository.NewNullTime(s.nowFn())
			if err := s.entryRepository.Update(ctx, entry); err != nil {
				return fmt.Errorf("error updating entry: %w", err)
			}
		}
	}
	return nil
}

func (s FeedService) Update(ctx context.Context, feed repository.Feed) error {
	feed.UpdatedAt = repository.NewNullTime(s.nowFn())
	if err := s.feedRepository.Update(ctx, feed); err != nil {
		return fmt.Errorf("error updating feed: %w", err)
	}
	return nil
}

func (s FeedService) Delete(ctx context.Context, feed repository.Feed) error {
	feed.DeletedAt = repository.NewNullTime(s.nowFn())
	if err := s.Update(ctx, feed); err != nil {
		return fmt.Errorf("error deleting feed: %w", err)
	}
	return nil
}

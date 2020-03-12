package webrss

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Alkemic/webrss/repository"

	"github.com/Alkemic/webrss/feed_fetcher"
)

type feedFetcher interface {
	Fetch(url string) (feed_fetcher.Feed, error)
}

type FeedService struct {
	nowFn           func() time.Time
	feedRepository  feedRepository
	entryRepository entryRepository
	feedFetcher     feedFetcher
}

func NewFeedService(feedRepository feedRepository, entryRepository entryRepository, feedFetcher feedFetcher) *FeedService {
	return &FeedService{
		nowFn:           time.Now,
		feedRepository:  feedRepository,
		entryRepository: entryRepository,
		feedFetcher:     feedFetcher,
	}
}

func (s FeedService) Create(feedURL string, categoryID int64) error {
	feeder, err := s.feedFetcher.Fetch(feedURL)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}
	feed := feeder.Feed()
	entries := feeder.Entries()
	now := repository.NewTime(s.nowFn())
	feed.CategoryID = categoryID
	feed.CreatedAt = now
	feed.LastReadAt = repository.NewTime(time.Date(1900, 1, 1, 1, 1, 1, 1, time.UTC))

	if err := s.feedRepository.Begin(); err != nil {
		return fmt.Errorf("cannot start transation when creating new feed: %w", err)
	}
	defer s.feedRepository.Rollback()

	feedID, err := s.feedRepository.Create(feed)
	if err != nil {
		return fmt.Errorf("error creating new feed: %w", err)
	}
	if err := s.SaveEntries(feedID, entries); err != nil {
		return fmt.Errorf("error saving entries: %w", err)
	}
	if err := s.feedRepository.Commit(); err != nil {
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

func (s FeedService) SaveEntries(feedID int64, entries []repository.Entry) error {
	now := repository.NewTime(s.nowFn())
	for _, entry := range entries {
		existingEntry, err := s.entryRepository.GetByURL(entry.Link)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("error fetching entry: %w", err)
		} else if err == sql.ErrNoRows {
			entry.FeedID = feedID
			entry.CreatedAt = now
			if err := s.entryRepository.Create(entry); err != nil {
				return fmt.Errorf("error creating entry: %w", err)
			}
		} else {
			entry = updateEntry(existingEntry, entry)
			if err := s.entryRepository.Update(entry); err != nil {
				return fmt.Errorf("error updating entry: %w", err)
			}
		}
	}
	return nil
}

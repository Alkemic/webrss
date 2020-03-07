package webrss

import (
	"fmt"
	"time"

	"github.com/Alkemic/webrss/repository"

	"github.com/Alkemic/webrss/feed_fetcher"
)

type feedFetcher interface {
	Fetch(url string) (feed_fetcher.Feed, error)
}

type FeedService struct {
	feedRepository  feedRepository
	entryRepository entryRepository
	feedFetcher     feedFetcher
}

func NewFeedService(feedRepository feedRepository, entryRepository entryRepository, feedFetcher feedFetcher) *FeedService {
	return &FeedService{
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
	now := repository.NewTime(time.Now())
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

	for _, entry := range entries {
		entry.FeedID = feedID
		entry.CreatedAt = now
		if err := s.entryRepository.Create(entry); err != nil {
			return fmt.Errorf("error saving entry: %w", err)
		}
	}
	if err := s.feedRepository.Commit(); err != nil {
		return fmt.Errorf("cannot commit transation when creating new feed: %w", err)
	}
	return nil
}

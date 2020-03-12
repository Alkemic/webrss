package webrss

import (
	"fmt"
	"time"

	"github.com/Alkemic/webrss/repository"
)

type entryRepository interface {
	Get(id int64) (repository.Entry, error)
	GetByURL(url string) (repository.Entry, error)
	ListForFeed(feedID, page int64) ([]repository.Entry, error)
	Update(entry repository.Entry) error
	Create(entry repository.Entry) error
}

type EntryService struct {
	entryRepository entryRepository
	feedRepository  feedRepository
}

func NewEntryService(entryRepository entryRepository, feedRepository feedRepository) *EntryService {
	return &EntryService{
		entryRepository: entryRepository,
		feedRepository:  feedRepository,
	}
}

func (s EntryService) ListForFeed(feedID, page int64) ([]repository.Entry, error) {
	entries, err := s.entryRepository.ListForFeed(feedID, page)
	if err != nil {
		return nil, fmt.Errorf("error fetching entries for feed %d: %w", feedID, err)
	}

	feed, err := s.feedRepository.Get(feedID)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch feed for entries: %w", err)
	}

	for i, _ := range entries {
		entries[i].Feed = feed
		entries[i].NewEntry = entries[i].CreatedAt.After(feed.LastReadAt)
	}

	feed.LastReadAt.Time = time.Now()
	if err := s.feedRepository.Update(feed); err != nil {
		return nil, fmt.Errorf("cannot fetch feed for entries: %w", err)
	}

	return entries, nil
}

func (s EntryService) Get(id int64) (repository.Entry, error) {
	entry, err := s.entryRepository.Get(id)
	if err != nil {
		return repository.Entry{}, fmt.Errorf("error getting entry: %w", err)
	}

	if !entry.ReadAt.Valid {
		entry.ReadAt.Valid = true
		entry.ReadAt.Time = time.Now()
		if err := s.entryRepository.Update(entry); err != nil {
			return entry, fmt.Errorf("error updating entry with 'read_at': %w", err)
		}
	}

	entry.Feed, err = s.feedRepository.Get(entry.FeedID)
	if err != nil {
		return entry, fmt.Errorf("cannot fetch feed for entry: %w", err)
	}

	entry.NewEntry = entry.CreatedAt.After(entry.Feed.LastReadAt)

	return entry, nil
}

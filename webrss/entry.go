package webrss

import (
	"context"
	"fmt"
	"time"

	"github.com/Alkemic/webrss/repository"
)

type entryRepository interface {
	Get(ctx context.Context, id int64) (repository.Entry, error)
	GetByURL(ctx context.Context, url string, feedID int64) (repository.Entry, error)
	ListForFeed(ctx context.Context, feedID, page int64, perPage int) ([]repository.Entry, error)
	ListForPhrase(ctx context.Context, phrase string, page int64, perPage int) ([]repository.Entry, error)
	Update(ctx context.Context, entry repository.Entry) error
	Create(ctx context.Context, entry repository.Entry) error
}

func (s WebRSSService) ListEntriesForFeed(ctx context.Context, feedID, page int64, perPage int) ([]repository.Entry, error) {
	entries, err := s.entryRepository.ListForFeed(ctx, feedID, page, perPage)
	if err != nil {
		return nil, fmt.Errorf("error fetching entries for feed %d: %w", feedID, err)
	}

	feed, err := s.feedRepository.Get(ctx, feedID)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch feed for entries: %w", err)
	}

	for i, _ := range entries {
		entries[i].Feed = feed
		entries[i].NewEntry = entries[i].CreatedAt.After(feed.LastReadAt)
	}

	feed.LastReadAt.Time = time.Now()
	if err := s.feedRepository.Update(ctx, feed); err != nil {
		return nil, fmt.Errorf("cannot fetch feed for entries: %w", err)
	}

	return entries, nil
}

func (s WebRSSService) GetEntry(ctx context.Context, id int64) (repository.Entry, error) {
	entry, err := s.entryRepository.Get(ctx, id)
	if err != nil {
		return repository.Entry{}, fmt.Errorf("error getting entry: %w", err)
	}

	if !entry.ReadAt.Valid {
		entry.ReadAt.Valid = true
		entry.ReadAt.Time = time.Now()
		if err := s.entryRepository.Update(ctx, entry); err != nil {
			return entry, fmt.Errorf("error updating entry with 'read_at': %w", err)
		}
	}

	entry.Feed, err = s.feedRepository.Get(ctx, entry.FeedID)
	if err != nil {
		return entry, fmt.Errorf("cannot fetch feed for entry: %w", err)
	}

	entry.NewEntry = entry.CreatedAt.After(entry.Feed.LastReadAt)

	return entry, nil
}

func (s WebRSSService) Search(ctx context.Context, phrase string, page int64, perPage int) ([]repository.Entry, error) {
	entries, err := s.entryRepository.ListForPhrase(ctx, phrase, page, perPage)
	if err != nil {
		return nil, fmt.Errorf("error fetching entries for phrase %s: %w", phrase, err)
	}

	return entries, nil
}

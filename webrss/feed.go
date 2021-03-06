package webrss

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Alkemic/webrss/favicon"
	"github.com/Alkemic/webrss/repository"
)

func (s WebRSSService) GetFeed(ctx context.Context, id int64) (repository.Feed, error) {
	feed, err := s.feedRepository.Get(ctx, id)
	if err != nil {
		return repository.Feed{}, fmt.Errorf("error getting feed: %w", err)
	}
	return feed, nil
}

func (s WebRSSService) CreateFeed(ctx context.Context, feedURL string, categoryID int64) error {
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

func (s WebRSSService) SaveEntries(ctx context.Context, feedID int64, entries []repository.Entry) error {
	now := repository.NewTime(s.nowFn())
	for _, entry := range entries {
		existingEntry, err := s.entryRepository.GetByURL(ctx, entry.Link, feedID)
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

func (s WebRSSService) UpdateFeed(ctx context.Context, feed repository.Feed) error {
	feeder, err := s.feedFetcher.Fetch(ctx, feed.FeedUrl)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}
	entries := feeder.Entries(ctx)
	log.Println("favicon url:", feeder.Feed(ctx).SiteFaviconUrl.String)

	if feed.SiteFaviconUrl.String != "" {
		faviconContent, err := favicon.DoRequest(ctx, s.httpClient, feed.SiteFaviconUrl.String)
		if err != nil {
			s.logger.Println("cannot download favicon: ", err)
		}
		feed.SiteFavicon = repository.NewNullString(base64.StdEncoding.EncodeToString(faviconContent))
	} else {
		feed.SiteFavicon.String = ""
	}

	if err := s.transactionRepository.Begin(ctx); err != nil {
		return fmt.Errorf("cannot start transation when creating new feed: %w", err)
	}
	defer s.transactionRepository.Rollback(ctx)

	feed.UpdatedAt = repository.NewNullTime(s.nowFn())
	if err := s.feedRepository.Update(ctx, feed); err != nil {
		return fmt.Errorf("error updating feed: %w", err)
	}
	if err := s.SaveEntries(ctx, feed.ID, entries); err != nil {
		return fmt.Errorf("error saving entries: %w", err)
	}
	if err := s.transactionRepository.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit transation when creating new feed: %w", err)
	}
	return nil
}

func (s WebRSSService) DeleteFeed(ctx context.Context, feed repository.Feed) error {
	now := repository.NewNullTime(s.nowFn())
	feed.UpdatedAt = now
	feed.DeletedAt = now
	if err := s.feedRepository.Update(ctx, feed); err != nil {
		return fmt.Errorf("error deleting feed: %w", err)
	}
	return nil
}

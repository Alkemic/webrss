package updater

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/sync/errgroup"

	"github.com/Alkemic/webrss/feed_fetcher"
	"github.com/Alkemic/webrss/repository"
)

type feedRepository interface {
	List(ctx context.Context) ([]repository.Feed, error)
}

type webrssService interface {
	SaveEntries(ctx context.Context, feedID int64, entries []repository.Entry) error
}

type feedFetcher interface {
	Fetch(ctx context.Context, url string) (feed_fetcher.Feed, error)
}

type UpdateService struct {
	feedRepository feedRepository
	webrssService  webrssService
	feedFetcher    feedFetcher
	logger         *log.Logger
}

func New(feedRepository feedRepository, webrssService webrssService, feedFetcher feedFetcher, logger *log.Logger) UpdateService {
	return UpdateService{
		feedRepository: feedRepository,
		webrssService:  webrssService,
		feedFetcher:    feedFetcher,
		logger:         logger,
	}
}

func (u UpdateService) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	feeds, err := u.feedRepository.List(ctx)
	if err != nil {
		return fmt.Errorf("cannot select feeds: %w", err)
	}
	for _, feed := range feeds {
		feed := feed
		g.Go(func() error {
			feeder, err := u.feedFetcher.Fetch(ctx, feed.FeedUrl)
			if err != nil {
				u.logger.Printf("error fetching feed %s: %v\n", feed.FeedUrl, err)
				return nil
			}
			entries := feeder.Entries(ctx)
			if err := u.webrssService.SaveEntries(ctx, feed.ID, entries); err != nil {
				return fmt.Errorf("cannot save entry: %w", err)
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return fmt.Errorf("got error during processing feeds: %w", err)
	}
	return nil
}

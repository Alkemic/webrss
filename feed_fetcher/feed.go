package feed_fetcher

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/Alkemic/webrss/favicon"
	"github.com/Alkemic/webrss/repository"
)

type Feed struct {
	parsedFeed *gofeed.Feed
	feedURL    string
	httpClient *http.Client
}

func New(parsedFeed *gofeed.Feed, httpClient *http.Client, feedURL string) Feed {
	return Feed{
		parsedFeed: parsedFeed,
		httpClient: httpClient,
		feedURL:    feedURL,
	}
}

func (f Feed) Feed(ctx context.Context) repository.Feed {
	faviconUrl, faviconContent, err := favicon.GetFavicon(ctx, f.httpClient, f.parsedFeed.Link)
	if errors.Is(err, favicon.ErrCannotParse) {
		faviconUrl, faviconContent, _ = favicon.GetFavicon(ctx, f.httpClient, f.parsedFeed.FeedLink)
	}
	if faviconContent.Valid {
		faviconContent.String = base64.StdEncoding.EncodeToString([]byte(faviconContent.String))
	}
	return repository.Feed{
		FeedUrl:        f.feedURL,
		FeedTitle:      f.parsedFeed.Title,
		FeedSubtitle:   repository.NewNullString(f.parsedFeed.Description),
		CreatedAt:      repository.NewTime(time.Now()),
		SiteFavicon:    faviconContent,
		SiteFaviconUrl: faviconUrl,
		SiteUrl:        repository.NewNullString(f.parsedFeed.Link),
	}
}

func (f Feed) Entries(_ context.Context) []repository.Entry {
	entries := make([]repository.Entry, 0, len(f.parsedFeed.Items))
	for _, item := range f.parsedFeed.Items {
		var publishedAt repository.Time
		if item.PublishedParsed != nil {
			publishedAt = repository.NewTime(*item.PublishedParsed)
		} else if item.UpdatedParsed != nil {
			publishedAt = repository.NewTime(*item.UpdatedParsed)
		}
		summary := item.Content
		if len(summary) == 0 {
			summary = item.Description
		}
		entries = append(entries, repository.Entry{
			Title:       item.Title,
			Summary:     repository.NewNullString(summary),
			Author:      parseAuthor(item.Author),
			Link:        item.Link,
			PublishedAt: publishedAt,
		})
	}
	return entries
}

func parseAuthor(feedAuthor *gofeed.Person) repository.NullString {
	var author repository.NullString
	if feedAuthor == nil || feedAuthor.Name == "" {
		return author
	}
	return repository.NewNullString(feedAuthor.Name)
}

package feed_fetcher

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mmcdole/gofeed"
)

//const defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"
const defaultUserAgent = "WebRSS parser (https://github.com/Alkemic/webrss)"

type FeedFetcher struct {
	parser     *gofeed.Parser
	httpClient *http.Client
}

func NewFeedParser(parser *gofeed.Parser, httpClient *http.Client) *FeedFetcher {
	return &FeedFetcher{
		httpClient: httpClient,
		parser:     parser,
	}
}

func (f FeedFetcher) Fetch(ctx context.Context, url string) (Feed, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Feed{}, fmt.Errorf("cannot create request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return Feed{}, fmt.Errorf("cannot execute request: %w", err)
	}
	defer resp.Body.Close()
	parsedFeed, err := f.parser.Parse(resp.Body)
	if err != nil {
		return Feed{}, fmt.Errorf("cannot parse feed data: %w", err)
	}

	return NewFeed(parsedFeed, f.httpClient, url), nil
}

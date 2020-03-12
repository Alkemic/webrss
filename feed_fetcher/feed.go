package feed_fetcher

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/mmcdole/gofeed"

	"github.com/Alkemic/webrss/repository"
)

type Feed struct {
	parsedFeed *gofeed.Feed
	feedURL    string
	httpClient *http.Client
}

func NewFeed(parsedFeed *gofeed.Feed, httpClient *http.Client, feedURL string) Feed {
	return Feed{
		parsedFeed: parsedFeed,
		httpClient: httpClient,
		feedURL:    feedURL,
	}
}

func (f Feed) Feed(ctx context.Context) repository.Feed {
	faviconUrl, favicon := f.getFavicon(ctx)
	return repository.Feed{
		FeedUrl:        f.feedURL,
		FeedTitle:      f.parsedFeed.Title,
		FeedSubtitle:   repository.NewNullString(f.parsedFeed.Description),
		CreatedAt:      repository.NewTime(time.Now()),
		SiteFavicon:    favicon,
		SiteFaviconUrl: faviconUrl,
		SiteUrl:        repository.NewNullString(f.parsedFeed.Link),
	}
}

func (f Feed) getFavicon(ctx context.Context) (repository.NullString, repository.NullString) {
	parsedUrl, err := url.Parse(f.parsedFeed.Link)
	if err != nil {
		parsedUrl, err = url.Parse(f.parsedFeed.FeedLink)
		if err != nil {
			return repository.NullString{}, repository.NullString{}
		}
	}

	baseUrl := parsedUrl.Scheme + "://" + parsedUrl.Host
	faviconUrl := baseUrl + "/favicon.ico"
	favicon, err := f.doRequest(ctx, faviconUrl)
	if err != nil {
		if faviconUrl, err = f.getFaviconURL(ctx, baseUrl); err != nil {
			return repository.NullString{}, repository.NullString{}
		}
		if favicon, err = f.doRequest(ctx, faviconUrl); err != nil {
			return repository.NullString{}, repository.NullString{}
		}
	}
	return repository.NewNullString(faviconUrl), repository.NewNullString(string(favicon))
}

func (f Feed) getFaviconURL(ctx context.Context, url string) (string, error) {
	body, err := f.doRequest(ctx, url)
	if err != nil {
		return "", fmt.Errorf("error doing request: %w", err)
	}
	doc, err := htmlquery.Parse(bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("error parsing html: %w", err)
	}
	faviconNode, err := htmlquery.Query(doc, `//link[@rel="icon" or @rel="shortcut icon"]/@href`)
	if err != nil {
		return "", fmt.Errorf("error quering xpath: %w", err)
	}
	return faviconNode.Data, nil
}

func (f Feed) doRequest(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	faviconRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}
	if faviconRaw == nil {
		return nil, fmt.Errorf("got empty body")
	}
	return faviconRaw, nil
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
		entries = append(entries, repository.Entry{
			Title:       item.Title,
			Summary:     repository.NewNullString(item.Content),
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

package favicon

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/antchfx/htmlquery"

	"github.com/Alkemic/webrss/repository"
)

var (
	ErrCannotParse = errors.New("cannot parse url")
)

//const defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"
const defaultUserAgent = "WebRSS parser (https://github.com/Alkemic/webrss)"

func GetFavicon(ctx context.Context, httpClient *http.Client, link string) (repository.NullString, repository.NullString, error) {
	parsedUrl, err := url.Parse(link)
	if err != nil {
		return repository.NullString{}, repository.NullString{}, fmt.Errorf("%s: %w", err.Error(), ErrCannotParse)
	}

	baseUrl := parsedUrl.Scheme + "://" + parsedUrl.Host
	faviconUrl := baseUrl + "/favicon.ico"
	favicon, err := DoRequest(ctx, httpClient, faviconUrl)
	if err != nil {
		if faviconUrl, err = getFaviconURL(ctx, httpClient, baseUrl); err != nil {
			return repository.NullString{}, repository.NullString{}, fmt.Errorf("cannot fetch favicon: %w", err)
		}
		if favicon, err = DoRequest(ctx, httpClient, faviconUrl); err != nil {
			return repository.NullString{}, repository.NullString{}, fmt.Errorf("cannot fetch site: %w", err)
		}
	}
	return repository.NewNullString(faviconUrl), repository.NewNullString(string(favicon)), nil
}

func getFaviconURL(ctx context.Context, httpClient *http.Client, url string) (string, error) {
	body, err := DoRequest(ctx, httpClient, url)
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

func DoRequest(ctx context.Context, httpClient *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := httpClient.Do(req)
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

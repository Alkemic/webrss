package webrss

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Alkemic/webrss/feed_fetcher"
	"github.com/Alkemic/webrss/repository"
)

type feedFetcher interface {
	Fetch(ctx context.Context, url string) (feed_fetcher.Feed, error)
}

type categoryRepository interface {
	List(ctx context.Context, params ...string) ([]repository.Category, error)
	Get(ctx context.Context, id int64) (repository.Category, error)
	Update(ctx context.Context, category repository.Category) error
	GetNextByOrder(ctx context.Context, order int) (repository.Category, error)
	GetPrevByOrder(ctx context.Context, order int) (repository.Category, error)
	Create(ctx context.Context, category repository.Category) error
	SelectMaxOrder(ctx context.Context) (int, error)
}

type feedRepository interface {
	Get(ctx context.Context, id int64) (repository.Feed, error)
	Create(ctx context.Context, feed repository.Feed) (int64, error)
	ListForCategories(ctx context.Context, ids []int64) ([]repository.Feed, error)
	List(ctx context.Context) ([]repository.Feed, error)
	Update(ctx context.Context, entry repository.Feed) error
}

type transactionRepository interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type WebRSSService struct {
	nowFn                 func() time.Time
	logger                *log.Logger
	feedRepository        feedRepository
	entryRepository       entryRepository
	categoryRepository    categoryRepository
	transactionRepository transactionRepository
	feedFetcher           feedFetcher
	httpClient            *http.Client
}

func NewService(
	logger *log.Logger,
	categoryRepository categoryRepository, feedRepository feedRepository,
	entryRepository entryRepository, transactionRepository transactionRepository,
	httpClient *http.Client, feedFetcher feedFetcher,
) *WebRSSService {
	return &WebRSSService{
		nowFn:                 time.Now,
		logger:                logger,
		feedRepository:        feedRepository,
		entryRepository:       entryRepository,
		categoryRepository:    categoryRepository,
		transactionRepository: transactionRepository,
		httpClient:            httpClient,
		feedFetcher:           feedFetcher,
	}
}

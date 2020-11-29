package webrss

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/Alkemic/webrss/feed_fetcher"
	"github.com/Alkemic/webrss/repository"
)

type feedFetcherMock struct{}

func (m feedFetcherMock) Fetch(ctx context.Context, url string) (feed_fetcher.Feed, error) {
	panic("implement me!")
}

type entryRepositoryMock struct {
	getEntryByURLResp map[string]repository.Entry
	getEntryByURLErr  map[string]error
	updateEntries     []repository.Entry
	updateErr         error
	createEntries     []repository.Entry
	createErr         error
}

func (m *entryRepositoryMock) Get(ctx context.Context, id int64) (repository.Entry, error) {
	panic("implement me!")
}

func (m *entryRepositoryMock) GetByURL(ctx context.Context, url string, feedID int64) (repository.Entry, error) {
	return m.getEntryByURLResp[url], m.getEntryByURLErr[url]
}

func (m *entryRepositoryMock) ListForFeed(ctx context.Context, feedID, page int64, perPage int) ([]repository.Entry, error) {
	panic("implement me!")
}

func (m *entryRepositoryMock) ListForPhrase(ctx context.Context, phrase string, page int64, perPage int) ([]repository.Entry, error) {
	panic("implement me!")
}

func (m *entryRepositoryMock) Update(ctx context.Context, entry repository.Entry) error {
	m.updateEntries = append(m.updateEntries, entry)
	return m.updateErr
}

func (m *entryRepositoryMock) Create(ctx context.Context, entry repository.Entry) error {
	m.createEntries = append(m.createEntries, entry)
	return m.createErr
}

type feedRepositoryMock struct{}

func (m *feedRepositoryMock) Get(ctx context.Context, id int64) (repository.Feed, error) {
	panic("implement me!")
}

func (m *feedRepositoryMock) List(ctx context.Context) ([]repository.Feed, error) {
	panic("implement me!")
}

func (m *feedRepositoryMock) Create(ctx context.Context, feed repository.Feed) (int64, error) {
	panic("implement me!")
}

func (m *feedRepositoryMock) ListForCategories(ctx context.Context, ids []int64) ([]repository.Feed, error) {
	panic("implement me!")
}

func (m *feedRepositoryMock) Update(ctx context.Context, entry repository.Feed) error {
	panic("implement me!")
}

type transactionRepositoryMock struct{}

func (m *transactionRepositoryMock) Begin(ctx context.Context) error {
	panic("implement me!")
}

func (m *transactionRepositoryMock) Commit(ctx context.Context) error {
	panic("implement me!")
}

func (m *transactionRepositoryMock) Rollback(ctx context.Context) error {
	panic("implement me!")
}

func TestFeedService_SaveEntries(t *testing.T) {
	type check func(error, *entryRepositoryMock, *testing.T)
	checks := func(cs ...check) []check { return cs }
	hasNoError := func(err error, mock *entryRepositoryMock, t *testing.T) {
		t.Helper()
		if err != nil {
			t.Errorf("Expected err to be nil, but got '%v'", err)
		}
	}
	hasUpdatedEntries := func(expectedResult []repository.Entry) check {
		return func(err error, mock *entryRepositoryMock, t *testing.T) {
			t.Helper()
			if !reflect.DeepEqual(mock.updateEntries, expectedResult) {
				t.Errorf("Expected result length to be '%d', but got '%d'", len(expectedResult), len(mock.updateEntries))
				t.Errorf("Expected result to be '%+v', but got '%+v'", expectedResult, mock.updateEntries)
			}
		}
	}
	hasCreatedEntries := func(expectedResult []repository.Entry) check {
		return func(err error, mock *entryRepositoryMock, t *testing.T) {
			t.Helper()
			if !reflect.DeepEqual(mock.createEntries, expectedResult) {
				t.Errorf("Expected result length to be '%d', but got '%d'", len(expectedResult), len(mock.createEntries))
				t.Errorf("Expected result to be '%+v', but got '%+v'", expectedResult, mock.createEntries)
			}
		}
	}
	hasError := func(expectedErr error) check {
		return func(err error, mock *entryRepositoryMock, t *testing.T) {
			t.Helper()
			if !errors.Is(err, expectedErr) {
				t.Errorf("Expected error to be '%v', but got '%v'", expectedErr, err)
			}
		}
	}
	hasErrorMsg := func(expectedErr string) check {
		return func(err error, mock *entryRepositoryMock, t *testing.T) {
			t.Helper()
			if err.Error() != expectedErr {
				t.Errorf("Expected error to be '%v', but got '%v'", expectedErr, err)
			}
		}
	}
	mockedErr := errors.New("mocked error")
	tests := []struct {
		name    string
		feedID  int64
		entries []repository.Entry
		// url => entry
		getEntryByURLResp map[string]repository.Entry
		getEntryByURLErr  map[string]error
		updateErr         error
		createErr         error
		ctx               context.Context

		checks []check
	}{{
		name:   "nothing on empty entries",
		checks: checks(hasNoError),
	}, {
		name:             "create new entries",
		feedID:           12,
		getEntryByURLErr: map[string]error{"link1": sql.ErrNoRows},
		entries: []repository.Entry{{
			Title:   "title 1",
			Author:  repository.NewNullString("author 1"),
			Summary: repository.NewNullString("summary 1"),
			Link:    "link1",
		}},
		checks: checks(
			hasNoError,
			hasCreatedEntries([]repository.Entry{{
				Title:     "title 1",
				Author:    repository.NewNullString("author 1"),
				Summary:   repository.NewNullString("summary 1"),
				Link:      "link1",
				CreatedAt: repository.NewTime(time.Date(2016, 3, 19, 7, 56, 35, 0, time.UTC)),
				FeedID:    12,
			}}),
		),
	}, {
		name:   "update entries",
		feedID: 13,
		getEntryByURLResp: map[string]repository.Entry{"link1": {
			Title:     "title 1",
			Author:    repository.NewNullString("author 1"),
			Summary:   repository.NewNullString("summary 1"),
			Link:      "link1",
			CreatedAt: repository.NewTime(time.Date(2014, 3, 19, 7, 56, 35, 0, time.UTC)),
			FeedID:    12,
		}},
		entries: []repository.Entry{{
			Title:   "new title",
			Author:  repository.NewNullString("new author"),
			Summary: repository.NewNullString("new summary"),
			Link:    "link1",
		}},
		checks: checks(
			hasNoError,
			hasUpdatedEntries([]repository.Entry{{
				Title:     "new title",
				Author:    repository.NewNullString("new author"),
				Summary:   repository.NewNullString("new summary"),
				Link:      "link1",
				CreatedAt: repository.NewTime(time.Date(2014, 3, 19, 7, 56, 35, 0, time.UTC)),
				UpdatedAt: repository.NewNullTime(time.Date(2016, 3, 19, 7, 56, 35, 0, time.UTC)),
				FeedID:    12,
			}}),
		),
	}, {
		name:             "update and create entry",
		feedID:           13,
		getEntryByURLErr: map[string]error{"link2": sql.ErrNoRows},
		getEntryByURLResp: map[string]repository.Entry{"link1": {
			Title:     "title 1",
			Author:    repository.NewNullString("author 1"),
			Summary:   repository.NewNullString("summary 1"),
			Link:      "link1",
			CreatedAt: repository.NewTime(time.Date(2014, 3, 19, 7, 56, 35, 0, time.UTC)),
			FeedID:    12,
		}},
		entries: []repository.Entry{{
			Title:   "new title",
			Author:  repository.NewNullString("new author"),
			Summary: repository.NewNullString("new summary"),
			Link:    "link1",
		}, {
			Title:   "title 2",
			Author:  repository.NewNullString("author 2"),
			Summary: repository.NewNullString("summary 2"),
			Link:    "link2",
		}},
		checks: checks(
			hasNoError,
			hasUpdatedEntries([]repository.Entry{{
				Title:     "new title",
				Author:    repository.NewNullString("new author"),
				Summary:   repository.NewNullString("new summary"),
				Link:      "link1",
				CreatedAt: repository.NewTime(time.Date(2014, 3, 19, 7, 56, 35, 0, time.UTC)),
				UpdatedAt: repository.NewNullTime(time.Date(2016, 3, 19, 7, 56, 35, 0, time.UTC)),
				FeedID:    12,
			}}),
			hasCreatedEntries([]repository.Entry{{
				Title:     "title 2",
				Author:    repository.NewNullString("author 2"),
				Summary:   repository.NewNullString("summary 2"),
				Link:      "link2",
				CreatedAt: repository.NewTime(time.Date(2016, 3, 19, 7, 56, 35, 0, time.UTC)),
				FeedID:    13,
			}}),
		),
	}, {
		name:             "error selecting entry",
		feedID:           13,
		getEntryByURLErr: map[string]error{"link1": mockedErr},
		entries: []repository.Entry{{
			Title:   "new title",
			Author:  repository.NewNullString("new author"),
			Summary: repository.NewNullString("new summary"),
			Link:    "link1",
		}},
		checks: checks(
			hasError(mockedErr),
			hasErrorMsg("error fetching entry: mocked error"),
		),
	}, {
		name:             "error while creating new entries",
		feedID:           12,
		getEntryByURLErr: map[string]error{"link1": sql.ErrNoRows},
		entries: []repository.Entry{{
			Title:   "title 1",
			Author:  repository.NewNullString("author 1"),
			Summary: repository.NewNullString("summary 1"),
			Link:    "link1",
		}},
		createErr: mockedErr,
		checks: checks(
			hasError(mockedErr),
			hasErrorMsg("error creating entry: mocked error"),
		),
	}, {
		name:   "error while updating entries",
		feedID: 13,
		getEntryByURLResp: map[string]repository.Entry{"link1": {
			Title:     "title 1",
			Author:    repository.NewNullString("author 1"),
			Summary:   repository.NewNullString("summary 1"),
			Link:      "link1",
			CreatedAt: repository.NewTime(time.Date(2014, 3, 19, 7, 56, 35, 0, time.UTC)),
			FeedID:    12,
		}},
		updateErr: mockedErr,
		entries: []repository.Entry{{
			Title:   "new title",
			Author:  repository.NewNullString("new author"),
			Summary: repository.NewNullString("new summary"),
			Link:    "link1",
		}},
		checks: checks(
			hasError(mockedErr),
			hasErrorMsg("error updating entry: mocked error"),
		),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockedEntryRepository := &entryRepositoryMock{
				getEntryByURLResp: tt.getEntryByURLResp,
				getEntryByURLErr:  tt.getEntryByURLErr,
				updateErr:         tt.updateErr,
				createErr:         tt.createErr,
			}
			mockedFeedRepository := &feedRepositoryMock{}
			s := WebRSSService{
				nowFn: func() time.Time {
					return time.Date(2016, 3, 19, 7, 56, 35, 0, time.UTC)
				},
				feedRepository:  mockedFeedRepository,
				entryRepository: mockedEntryRepository,
				feedFetcher:     feedFetcherMock{},
			}
			err := s.SaveEntries(tt.ctx, tt.feedID, tt.entries)
			for _, ch := range tt.checks {
				ch(err, mockedEntryRepository, t)
			}
		})
	}
}

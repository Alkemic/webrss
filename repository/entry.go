package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	selectEntriesForFeedQuery = `
SELECT *
FROM entry e
where e.deleted_at is null and e.feed_id = ?
ORDER BY e.published_at DESC
LIMIT ? OFFSET ?;`
	getEntryQuery      = `SELECT * FROM entry e where e.deleted_at is null and id = ?;`
	getEntryByURLQuery = `SELECT * FROM entry e where e.deleted_at is null and link = ?;`
	updateEntryQuery   = `
update entry 
set title = :title, author = :author, summary = :summary, link = :link, published_at = :published_at, 
feed_id = :feed_id, read_at = :read_at, created_at = :created_at, updated_at = :updated_at, deleted_at = :deleted_at 
where id = :id and deleted_at is null;`
	createEntryQuery = `insert into entry(title, author, summary, link, published_at, feed_id, read_at, created_at)
values (:title, :author, :summary, :link, :published_at, :feed_id, :read_at, :created_at);`
)

type entryRepository struct {
	db      *sqlx.DB
	perPage int64
}

func NewEntryRepository(db *sqlx.DB, perPage int64) *entryRepository {
	return &entryRepository{
		db:      db,
		perPage: perPage,
	}
}

func (r *entryRepository) Get(ctx context.Context, id int64) (Entry, error) {
	entry := Entry{}
	if err := r.db.GetContext(ctx, &entry, getEntryQuery, id); err != nil {
		return Entry{}, fmt.Errorf("cannot fetch entry (id=%d): %w", id, err)
	}
	return entry, nil
}

func (r *entryRepository) GetByURL(ctx context.Context, url string) (Entry, error) {
	entry := Entry{}
	if err := r.db.GetContext(ctx, &entry, getEntryByURLQuery, url); err != nil {
		return Entry{}, fmt.Errorf("cannot fetch entry (url=%s): %w", url, err)
	}
	return entry, nil
}

func (r *entryRepository) ListForFeed(ctx context.Context, feedID, page int64) ([]Entry, error) {
	entries := []Entry{}
	if err := r.db.SelectContext(ctx, &entries, selectEntriesForFeedQuery, feedID, r.perPage, r.perPage*(page-1)); err != nil {
		return nil, fmt.Errorf("cannot select entries: %w", err)
	}
	return entries, nil
}

func (r *entryRepository) Create(ctx context.Context, entry Entry) error {
	if _, err := r.db.NamedExecContext(ctx, createEntryQuery, entry); err != nil {
		return fmt.Errorf("cannot create entry: %w", err)
	}
	return nil
}

func (r *entryRepository) Update(ctx context.Context, entry Entry) error {
	if _, err := r.db.NamedExecContext(ctx, updateEntryQuery, entry); err != nil {
		return fmt.Errorf("cannot update entry: %w", err)
	}
	return nil
}

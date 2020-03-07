package repository

import (
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
	getEntryQuery    = `SELECT * FROM entry e where e.deleted_at is null and id = ?;`
	updateEntryQuery = `
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

func (r *entryRepository) Get(id int64) (Entry, error) {
	entry := Entry{}
	if err := r.db.Get(&entry, getEntryQuery, id); err != nil {
		return Entry{}, fmt.Errorf("cannot fetch entry: %w", err)
	}
	return entry, nil
}
func (r *entryRepository) ListForFeed(feedID, page int64) ([]Entry, error) {
	entries := []Entry{}
	if err := r.db.Select(&entries, selectEntriesForFeedQuery, feedID, r.perPage, r.perPage*(page-1)); err != nil {
		return nil, fmt.Errorf("cannot select entries: %w", err)
	}
	return entries, nil
}

func (r *entryRepository) Create(entry Entry) error {
	if _, err := r.db.NamedExec(createEntryQuery, entry); err != nil {
		return fmt.Errorf("cannot create entry: %w", err)
	}
	return nil
}

func (r *entryRepository) Update(entry Entry) error {
	if _, err := r.db.NamedExec(updateEntryQuery, entry); err != nil {
		return fmt.Errorf("cannot update entry: %w", err)
	}
	return nil
}

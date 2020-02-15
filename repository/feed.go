package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type feedRepository struct {
	db *sqlx.DB
}

func NewFeedRepository(db *sqlx.DB) *feedRepository {
	return &feedRepository{
		db: db,
	}
}

var _ = `
select
    f.id, (
        select max(e.published_at)
        from entry e
        where e.feed_id = f.id
    ) entry_max_published_at
from feed f
where f.deleted_at is null;
`

var (
	selectFeedsForCategoriesQuery = `
SELECT 
	*,
	(select count(*) from entry e where e.read_at is null and e.feed_id = f.id) un_read,
	f.last_read_at < (select max(e.published_at) from entry e where e.feed_id = f.id) new_entries
FROM feed f 
where f.deleted_at is null and f.category_id in (?) 
ORDER BY "f.order" ASC;`
)

func (r *feedRepository) ListForCategories(categoriesIDs []int64) ([]Feed, error) {
	query, args, err := sqlx.In(selectFeedsForCategoriesQuery, categoriesIDs)
	feeds := []Feed{}
	if err = r.db.Select(&feeds, r.db.Rebind(query), args...); err != nil {
		return nil, fmt.Errorf("cannot select feeds: %w", err)
	}
	return feeds, nil
}

func (r *feedRepository) Get(id int64) (Feed, error) {
	return Feed{}, nil
}

func (r *feedRepository) Delete(id int64) error {
	return nil
}

func (r *feedRepository) Update(category Feed) error {
	return nil
}

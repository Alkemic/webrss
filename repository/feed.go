package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	selectFeedsForCategoriesQuery = `
SELECT 
	*,
	(select count(*) from entry e where e.read_at is null and e.feed_id = f.id) un_read,
	coalesce(f.last_read_at < (select max(e.published_at) from entry e where e.feed_id = f.id), false) new_entries
FROM feed f 
where f.deleted_at is null and f.category_id in (?) 
ORDER BY "f.order" ASC;`
	selectFeedsQuery = `SELECT * FROM feed where deleted_at is null ORDER BY "order" ASC;`
	getFeedQuery     = `select * from feed where id = ? and deleted_at is null;`
	updateFeedQuery  = `
update feed 
set feed_title = :feed_title, feed_url = :feed_url, feed_image = :feed_image, feed_subtitle = :feed_subtitle, site_url = :site_url, 
site_favicon_url = :site_favicon_url, site_favicon = :site_favicon, category_id = :category_id, last_read_at = :last_read_at, 
created_at = :created_at, updated_at = :updated_at, deleted_at = :deleted_at
where id = :id and deleted_at is null;`
	createFeedQuery = `insert into feed (feed_title, feed_url, feed_image, feed_subtitle, site_url, site_favicon_url, site_favicon, category_id, last_read_at, created_at)
values (:feed_title, :feed_url, :feed_image, :feed_subtitle, :site_url, :site_favicon_url, :site_favicon, :category_id, :last_read_at, :created_at);`
)

type feedRepository struct {
	db *sqlx.DB
}

func NewFeedRepository(db *sqlx.DB) *feedRepository {
	return &feedRepository{
		db: db,
	}
}

func (r *feedRepository) Get(ctx context.Context, id int64) (Feed, error) {
	feed := Feed{}
	if err := r.db.GetContext(ctx, &feed, getFeedQuery, id); err != nil {
		return Feed{}, fmt.Errorf("cannot fetch feed: %w", err)
	}
	return feed, nil
}

func (r *feedRepository) ListForCategories(ctx context.Context, categoriesIDs []int64) ([]Feed, error) {
	query, args, err := sqlx.In(selectFeedsForCategoriesQuery, categoriesIDs)
	feeds := []Feed{}
	if err = r.db.SelectContext(ctx, &feeds, r.db.Rebind(query), args...); err != nil {
		return nil, fmt.Errorf("cannot select feeds: %w", err)
	}
	return feeds, nil
}

func (r *feedRepository) List(ctx context.Context) ([]Feed, error) {
	feeds := []Feed{}
	if err := r.db.SelectContext(ctx, &feeds, selectFeedsQuery); err != nil {
		return nil, fmt.Errorf("cannot select feeds: %w", err)
	}
	return feeds, nil
}

func (r *feedRepository) Update(ctx context.Context, feed Feed) error {
	if _, err := r.db.NamedExecContext(ctx, updateFeedQuery, feed); err != nil {
		return fmt.Errorf("cannot update feed: %w", err)
	}
	return nil
}

func (r *feedRepository) Create(ctx context.Context, feed Feed) (int64, error) {
	res, err := r.db.NamedExecContext(ctx, createFeedQuery, feed)
	if err != nil {
		return 0, fmt.Errorf("cannot create feed: %w", err)
	}
	lastInsertedID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("cannot fetch last inserted id: %w", err)
	}
	return lastInsertedID, nil
}

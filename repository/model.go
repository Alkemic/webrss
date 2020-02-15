package repository

import (
	"database/sql"
	"time"
)

type Category struct {
	ID        int64        `db:"id" json:"id"`        // | int(11)      | NO
	Title     string       `db:"title" json:"title"`  // | varchar(255) | NO
	Order     int          `db:"order" json:"-"`      // | int(11)      | NO
	CreatedAt time.Time    `db:"created_at" json:"-"` // | datetime     | NO
	UpdatedAt sql.NullTime `db:"updated_at" json:"-"` // | datetime     | YES
	DeletedAt sql.NullTime `db:"deleted_at" json:"-"` // | datetime     | YES

	Feeds []Feed `db:"-" json:"feeds"`
}

type Feed struct {
	ID             int64        `db:"id" json:"id"`                             // | int(11)      | NO   ||
	FeedTitle      string       `db:"feed_title" json:"feed_title"`             // | varchar(255) | NO  |
	FeedUrl        string       `db:"feed_url" json:"feed_url"`                 // | varchar(255) | NO  |
	FeedImage      NullString   `db:"feed_image" json:"feed_image"`             // | varchar(255) | YES |
	FeedSubtitle   NullString   `db:"feed_subtitle" json:"feed_subtitle"`       // | longtext     | YES |
	SiteUrl        NullString   `db:"site_url" json:"site_url"`                 // | varchar(255) | YES |
	SiteFaviconUrl NullString   `db:"site_favicon_url" json:"site_favicon_url"` // | varchar(255) | YES |
	SiteFavicon    NullString   `db:"site_favicon" json:"site_favicon"`         // | text         | YES |
	CategoryID     int64        `db:"category_id" json:"category_id"`           // | int(11)      | YES  ||
	LastReadAt     time.Time    `db:"last_read_at" json:"-"`                    // | datetime     | NO  |
	CreatedAt      time.Time    `db:"created_at" json:"-"`                      // | datetime     | NO  |
	UpdatedAt      sql.NullTime `db:"updated_at" json:"-"`                      // | datetime     | YES |
	DeletedAt      sql.NullTime `db:"deleted_at" json:"-"`                      // | datetime     | YES  ||

	UnRead     int64 `db:"un_read" json:"un_read"`
	NewEntries int64 `db:"new_entries" json:"new_entries"`
}

type Entry struct {
	ID int64
}

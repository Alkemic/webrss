package repository

type Category struct {
	ID        int64    `db:"id" json:"id"`
	Title     string   `db:"title" json:"title"`
	Order     int      `db:"order" json:"-"`
	CreatedAt Time     `db:"created_at" json:"-"`
	UpdatedAt NullTime `db:"updated_at" json:"-"`
	DeletedAt NullTime `db:"deleted_at" json:"-"`

	Feeds []Feed `db:"-" json:"feeds"`
}

type Feed struct {
	ID             int64      `db:"id" json:"id"`
	FeedTitle      string     `db:"feed_title" json:"feed_title"`
	FeedUrl        string     `db:"feed_url" json:"feed_url"`
	FeedImage      NullString `db:"feed_image" json:"feed_image"`
	FeedSubtitle   NullString `db:"feed_subtitle" json:"feed_subtitle"`
	SiteUrl        NullString `db:"site_url" json:"site_url"`
	SiteFaviconUrl NullString `db:"site_favicon_url" json:"site_favicon_url"`
	SiteFavicon    NullString `db:"site_favicon" json:"site_favicon"`
	CategoryID     int64      `db:"category_id" json:"category_id"`
	LastReadAt     Time       `db:"last_read_at" json:"-"`
	CreatedAt      Time       `db:"created_at" json:"-"`
	UpdatedAt      NullTime   `db:"updated_at" json:"-"`
	DeletedAt      NullTime   `db:"deleted_at" json:"-"`

	UnRead     int64 `db:"un_read" json:"un_read"`
	NewEntries int64 `db:"new_entries" json:"new_entries"`
}

type Entry struct {
	ID          int64      `db:"id" json:"id"`
	Title       string     `db:"title" json:"title"`
	Author      NullString `db:"author" json:"author"`
	Summary     NullString `db:"summary" json:"summary"`
	Link        string     `db:"link" json:"link"`
	PublishedAt Time       `db:"published_at" json:"published_at"`
	FeedID      int64      `db:"feed_id" json:"feed_id"`
	ReadAt      NullTime   `db:"read_at" json:"read_at"`
	CreatedAt   Time       `db:"created_at" json:"-"`
	UpdatedAt   NullTime   `db:"updated_at" json:"-"`
	DeletedAt   NullTime   `db:"deleted_at" json:"-"`

	Feed     Feed `db:"-" json:"feed"`
	NewEntry bool `db:"-" json:"new_entry"`
}

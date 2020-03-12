package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	selectCategoriesQuery = `select * from category where deleted_at is null order by "order" asc;`
	maxCategoryOrderQuery = `select max(c.order) from category c where deleted_at is null;`
	createCategoryQuery   = "insert into category(title, `order`, created_at) values (:title, :order, :created_at);"
)

type categoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *categoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) List(ctx context.Context, params ...string) ([]Category, error) {
	categories := []Category{}
	err := r.db.SelectContext(ctx, &categories, selectCategoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("cannot select categories: %w", err)
	}
	return categories, nil
}

func (r *categoryRepository) Create(ctx context.Context, category Category) error {
	if _, err := r.db.NamedExecContext(ctx, createCategoryQuery, category); err != nil {
		return fmt.Errorf("cannot create category: %w", err)
	}
	return nil
}

func (r *categoryRepository) SelectMaxOrder(ctx context.Context) (int, error) {
	var maxOrder int
	if err := r.db.GetContext(ctx, &maxOrder, maxCategoryOrderQuery); err != nil {
		return -1, fmt.Errorf("cannot create category: %w", err)
	}
	return maxOrder, nil
}

func (r *categoryRepository) Get(ctx context.Context, id int64) (Category, error) {
	return Category{}, nil
}

func (r *categoryRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

func (r *categoryRepository) Update(ctx context.Context, category Category) error {
	return nil
}

func (r *categoryRepository) GetNextByOrder(ctx context.Context, category Category) (Category, error) {
	return Category{}, nil
}

func (r *categoryRepository) GetPrevByOrder(ctx context.Context, category Category) (Category, error) {
	return Category{}, nil
}

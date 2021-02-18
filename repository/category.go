package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	selectCategoriesQuery = "select * from category where deleted_at is null order by `order` asc;"
	selectCategoryQuery   = `select * from category where deleted_at is null and id = ?;`
	maxCategoryOrderQuery = `select coalesce(max(c.order), 0) max_order from category c where deleted_at is null;`
	createCategoryQuery   = "insert into category(title, `order`, created_at) values (:title, :order, :created_at);"
	getPrevQuery          = "select * from category where deleted_at is null and `order` < ? order by `order` desc limit 1;"
	getNextQuery          = "select * from category where deleted_at is null and `order` > ? order by `order` asc limit 1;"
	updateCategoryQuery   = "update category set title = :title, `order` = :order, created_at = :created_at, updated_at = :updated_at, deleted_at = :deleted_at where id = :id;"
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
	var category Category
	if err := r.db.GetContext(ctx, &category, selectCategoryQuery, id); err != nil {
		return Category{}, fmt.Errorf("cannot select category: %w", err)
	}
	return category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category Category) error {
	if _, err := r.db.NamedExecContext(ctx, updateCategoryQuery, category); err != nil {
		return fmt.Errorf("cannot update category: %w", err)
	}
	return nil
}

func (r *categoryRepository) GetNextByOrder(ctx context.Context, order int) (Category, error) {
	var category Category
	if err := r.db.GetContext(ctx, &category, getNextQuery, order); err != nil {
		return Category{}, fmt.Errorf("cannot select next category: %w", err)
	}
	return category, nil
}

func (r *categoryRepository) GetPrevByOrder(ctx context.Context, order int) (Category, error) {
	var category Category
	if err := r.db.GetContext(ctx, &category, getPrevQuery, order); err != nil {
		return Category{}, fmt.Errorf("cannot select previous category: %w", err)
	}
	return category, nil
}

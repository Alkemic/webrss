package repository

import (
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

func (r *categoryRepository) List(params ...string) ([]Category, error) {
	categories := []Category{}
	err := r.db.Select(&categories, selectCategoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("cannot select categories: %w", err)
	}
	return categories, nil
}

func (r *categoryRepository) Create(category Category) error {
	if _, err := r.db.NamedExec(createCategoryQuery, category); err != nil {
		return fmt.Errorf("cannot create category: %w", err)
	}
	return nil
}

func (r *categoryRepository) SelectMaxOrder() (int, error) {
	var maxOrder int
	if err := r.db.Get(&maxOrder, maxCategoryOrderQuery); err != nil {
		return -1, fmt.Errorf("cannot create category: %w", err)
	}
	return maxOrder, nil
}

func (r *categoryRepository) Get(id int64) (Category, error) {
	return Category{}, nil
}

func (r *categoryRepository) Delete(id int64) error {
	return nil
}

func (r *categoryRepository) Update(category Category) error {
	return nil
}

func (r *categoryRepository) GetNextByOrder(category Category) (Category, error) {
	return Category{}, nil
}

func (r *categoryRepository) GetPrevByOrder(category Category) (Category, error) {
	return Category{}, nil
}

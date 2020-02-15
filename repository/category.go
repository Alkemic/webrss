package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
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
	err := r.db.Select(&categories, `SELECT * FROM category ORDER BY "order" ASC;`)
	if err != nil {
		return nil, fmt.Errorf("cannot select categories: %w", err)
	}
	return categories, nil
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

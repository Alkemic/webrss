package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	getUserByIDQuery    = `select id, name, email, password from user where id = ?`
	getUserByEmailQuery = `select id, name, email, password from user where email = ?`
	updateUserQuery     = `update user set name = :name, email = :email, password = :password where id = :id;`
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r userRepository) GetByID(ctx context.Context, id int) (User, error) {
	user := User{}
	if err := r.db.GetContext(ctx, &user, getUserByIDQuery, id); err != nil {
		return User{}, fmt.Errorf("cannot fetch feed: %w", err)
	}
	return user, nil
}

func (r userRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	user := User{}
	if err := r.db.GetContext(ctx, &user, getUserByEmailQuery, email); err != nil {
		return User{}, fmt.Errorf("cannot fetch feed: %w", err)
	}
	return user, nil
}

func (r userRepository) Update(ctx context.Context, user User) error {
	if _, err := r.db.NamedExecContext(ctx, updateUserQuery, user); err != nil {
		return fmt.Errorf("cannot update user: %w", err)
	}
	return nil
}

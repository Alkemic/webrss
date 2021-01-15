package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usernameKey = "username"
	passwordKey = "password"

	getSettingsQuery = "select value from settings where `key` = ?;"
	setSettingsQuery = "update settings set value = ? where `key` = ?;"
)

type settingsRepository struct {
	db *sqlx.DB
}

func NewSettingsRepository(db *sqlx.DB) *settingsRepository {
	return &settingsRepository{db: db}
}

func (r *settingsRepository) Set(ctx context.Context, key, value string) error {
	if _, err := r.db.ExecContext(ctx, setSettingsQuery, value, key); err != nil {
		return fmt.Errorf("cannot fetch value of key %s: %w", key, err)
	}
	return nil
}

func (r *settingsRepository) SetTx(ctx context.Context, tx *sqlx.Tx, key, value string) error {
	if _, err := tx.ExecContext(ctx, setSettingsQuery, value, key); err != nil {
		return fmt.Errorf("cannot fetch value of key %s: %w", key, err)
	}
	return nil
}

func (r *settingsRepository) Get(ctx context.Context, key string) (string, error) {
	var value string
	if err := r.db.GetContext(ctx, &value, getSettingsQuery, key); err != nil {
		return "", fmt.Errorf("cannot fetch value of key %s: %w", key, err)
	}
	return value, nil
}

func (r *settingsRepository) GetUser(ctx context.Context) (User, error) {
	username, err := r.Get(ctx, usernameKey)
	if err != nil {
		return User{}, fmt.Errorf("cannot fetch username value: %s", err)
	}
	password, err := r.Get(ctx, passwordKey)
	if err != nil {
		return User{}, fmt.Errorf("cannot fetch password value: %s", err)
	}
	return User{
		Name:     username,
		Password: []byte(password),
	}, nil
}

func (r *settingsRepository) SaveUser(ctx context.Context, user User) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot begin transction: %w", err)
	}
	defer tx.Commit()
	if err := r.SetTx(ctx, tx, usernameKey, user.Name); err != nil {
		tx.Rollback()
		return fmt.Errorf("cannot update username: %w")
	}
	if err := r.SetTx(ctx, tx, passwordKey, string(user.Password)); err != nil {
		tx.Rollback()
		return fmt.Errorf("cannot update username: %w")
	}
	return nil
}

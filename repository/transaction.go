package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type transactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *transactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) Begin(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, "BEGIN;"); err != nil {
		return fmt.Errorf("cannot start transation: %w", err)
	}
	return nil
}

func (r *transactionRepository) Commit(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, "COMMIT;"); err != nil {
		return fmt.Errorf("cannot start transation: %w", err)
	}
	return nil
}

func (r *transactionRepository) Rollback(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, "ROLLBACK;"); err != nil {
		return fmt.Errorf("cannot start transation: %w", err)
	}
	return nil
}

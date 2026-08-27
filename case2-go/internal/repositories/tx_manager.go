package repositories

import (
	"context"
	"database/sql"
	"fmt"
)

// TxManager defines the contract for running database transactions.
type TxManager interface {
	WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error
}

type sqlTxManager struct {
	db *sql.DB
}

// NewTxManager creates a new SQL transaction manager.
func NewTxManager(db *sql.DB) TxManager {
	return &sqlTxManager{db: db}
}

// WithTransaction runs the provided function within an ACID transaction.
func (m *sqlTxManager) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after rollback
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

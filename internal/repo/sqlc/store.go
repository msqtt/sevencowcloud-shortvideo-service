package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions.
type Store interface {
	Querier
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store.
func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within database transactions.
// Rollback if function returns err, or Commit
// returns an error , if any.
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if err2 := tx.Rollback(); err != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, err2)
		}
	}
	return tx.Commit()
}

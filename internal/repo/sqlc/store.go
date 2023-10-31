package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Store provides all functions to execute db queries and transactions.
type Store interface {
	Querier
	CreateUserTx(ctx context.Context, param CreateUserTxParams) (CreateUserTxResult, error)
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

type CreateUserTxParams struct {
	NickName string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserTxResult struct {
	User User
	Profile Profile
}

func (store *SQLStore) CreateUserTx(ctx context.Context, param CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// add a profile
		r, err := q.AddProfile(ctx, AddProfileParams{
			AvatarLink: sql.NullString{String: "img/avatar/default.png", Valid: true},
		})
		if err != nil {
			return err
		}
		i, err := r.LastInsertId()
		if err != nil {
			return err
		}
		// create user
		params := AddUserParams{
			Nickname:  param.NickName,
			Email:     param.Email,
			Password:  param.Password,
			ProfileID: i,
		}
		r1, err := q.AddUser(ctx, params)
		if err != nil {
			return err
		}
		i2, err := r1.LastInsertId()	
		if err != nil {
			return err
		}

		result.User = User{
			ID: i2,
			Nickname: params.Nickname,	
			Email: param.Email,
			CreatedAt: time.Now(),
		}
		return nil
	})
	return result, err
}

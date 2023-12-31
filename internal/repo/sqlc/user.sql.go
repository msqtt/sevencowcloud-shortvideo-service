// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const addUser = `-- name: AddUser :execresult
INSERT INTO users (
	nickname, email, password, profile_id
) VALUES ( 
	?, ?, ?, ?
)
`

type AddUserParams struct {
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	ProfileID int64  `json:"profile_id"`
}

func (q *Queries) AddUser(ctx context.Context, arg AddUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, addUser,
		arg.Nickname,
		arg.Email,
		arg.Password,
		arg.ProfileID,
	)
}

const deleteUser = `-- name: DeleteUser :exec
UPDATE users
	SET is_deleted = 1 
	WHERE id = ?
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, nickname, email, password, profile_id, created_at, is_deleted FROM users
WHERE email = ? 
AND is_deleted = 0
ORDER By created_at DESC
LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Nickname,
		&i.Email,
		&i.Password,
		&i.ProfileID,
		&i.CreatedAt,
		&i.IsDeleted,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, nickname, email, password, profile_id, created_at, is_deleted FROM users
WHERE id = ? 
AND is_deleted = 0
LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Nickname,
		&i.Email,
		&i.Password,
		&i.ProfileID,
		&i.CreatedAt,
		&i.IsDeleted,
	)
	return i, err
}

const getUserByNickName = `-- name: GetUserByNickName :one
SELECT id, nickname, email, password, profile_id, created_at, is_deleted FROM users
WHERE nickname = ? 
AND is_deleted = 0
LIMIT 1
`

func (q *Queries) GetUserByNickName(ctx context.Context, nickname string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByNickName, nickname)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Nickname,
		&i.Email,
		&i.Password,
		&i.ProfileID,
		&i.CreatedAt,
		&i.IsDeleted,
	)
	return i, err
}

const updateNickName = `-- name: UpdateNickName :execresult
UPDATE users
	SET nickname = ?
	WHERE id =?
`

type UpdateNickNameParams struct {
	Nickname string `json:"nickname"`
	ID       int64  `json:"id"`
}

func (q *Queries) UpdateNickName(ctx context.Context, arg UpdateNickNameParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateNickName, arg.Nickname, arg.ID)
}

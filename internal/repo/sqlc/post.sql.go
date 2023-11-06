// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: post.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const addPost = `-- name: AddPost :execresult
INSERT INTO posts (
	title, description, user_id, video_id, created_at
) VALUES ( 
	?, ?, ?, ?, now()
)
`

type AddPostParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      int64  `json:"user_id"`
	VideoID     int64  `json:"video_id"`
}

func (q *Queries) AddPost(ctx context.Context, arg AddPostParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, addPost,
		arg.Title,
		arg.Description,
		arg.UserID,
		arg.VideoID,
	)
}

const deletePost = `-- name: DeletePost :exec
UPDATE posts
	SET is_deleted = 1
	WHERE id = ?
`

func (q *Queries) DeletePost(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deletePost, id)
	return err
}

const getPostByID = `-- name: GetPostByID :one
SELECT id, title, description, user_id, video_id, is_deleted, updated_at, created_at FROM posts
WHERE id = ?
AND is_deleted = 0
LIMIT 1
`

func (q *Queries) GetPostByID(ctx context.Context, id int64) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostByID, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.UserID,
		&i.VideoID,
		&i.IsDeleted,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getPostByUserID = `-- name: GetPostByUserID :many
SELECT id, title, description, user_id, video_id, is_deleted, updated_at, created_at FROM posts
WHERE user_id = ?
AND is_deleted = 0
LIMIT ?,?
`

type GetPostByUserIDParams struct {
	UserID int64 `json:"user_id"`
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

func (q *Queries) GetPostByUserID(ctx context.Context, arg GetPostByUserIDParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostByUserID, arg.UserID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Post{}
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.UserID,
			&i.VideoID,
			&i.IsDeleted,
			&i.UpdatedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchPostByTitle = `-- name: SearchPostByTitle :many
SELECT id, title, description, user_id, video_id, is_deleted, updated_at, created_at FROM posts
WHERE title like "%?%"
AND is_deleted = 0
LIMIT ?,?
`

type SearchPostByTitleParams struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

func (q *Queries) SearchPostByTitle(ctx context.Context, arg SearchPostByTitleParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, searchPostByTitle, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Post{}
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.UserID,
			&i.VideoID,
			&i.IsDeleted,
			&i.UpdatedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const testGetAll = `-- name: TestGetAll :many
SELECT id, title, description, user_id, video_id, is_deleted, updated_at, created_at, (SELECT count(id) FROM posts) total_size
FROM posts
WHERE is_deleted = 0
LIMIT ?, ?
`

type TestGetAllParams struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

type TestGetAllRow struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      int64     `json:"user_id"`
	VideoID     int64     `json:"video_id"`
	IsDeleted   int32     `json:"is_deleted"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
	TotalSize   int64     `json:"total_size"`
}

func (q *Queries) TestGetAll(ctx context.Context, arg TestGetAllParams) ([]TestGetAllRow, error) {
	rows, err := q.db.QueryContext(ctx, testGetAll, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []TestGetAllRow{}
	for rows.Next() {
		var i TestGetAllRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.UserID,
			&i.VideoID,
			&i.IsDeleted,
			&i.UpdatedAt,
			&i.CreatedAt,
			&i.TotalSize,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePostInfo = `-- name: UpdatePostInfo :exec
UPDATE posts
	SET title = ?, description = ?, updated_at = now()
	WHERE id = ?
`

type UpdatePostInfoParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ID          int64  `json:"id"`
}

func (q *Queries) UpdatePostInfo(ctx context.Context, arg UpdatePostInfoParams) error {
	_, err := q.db.ExecContext(ctx, updatePostInfo, arg.Title, arg.Description, arg.ID)
	return err
}
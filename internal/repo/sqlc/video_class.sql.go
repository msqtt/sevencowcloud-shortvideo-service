// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: video_class.sql

package db

import (
	"context"
)

const getAllVideoClass = `-- name: GetAllVideoClass :many
SELECT id, name, description, is_enabled FROM video_class
WHERE is_enabled = 1
`

func (q *Queries) GetAllVideoClass(ctx context.Context) ([]VideoClass, error) {
	rows, err := q.db.QueryContext(ctx, getAllVideoClass)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []VideoClass{}
	for rows.Next() {
		var i VideoClass
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.IsEnabled,
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

const getVideoClassByID = `-- name: GetVideoClassByID :one
SELECT id, name, description, is_enabled FROM video_class
WHERE id = ?
AND is_enabled = 1
`

func (q *Queries) GetVideoClassByID(ctx context.Context, id int32) (VideoClass, error) {
	row := q.db.QueryRowContext(ctx, getVideoClassByID, id)
	var i VideoClass
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.IsEnabled,
	)
	return i, err
}

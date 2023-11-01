// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.0
// source: user_captcha.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const addCaptcha = `-- name: AddCaptcha :execresult
INSERT INTO captcha (
	email, captcha, expired_at
) VALUES ( ?, ?, ? )
`

type AddCaptchaParams struct {
	Email     string    `json:"email"`
	Captcha   string    `json:"captcha"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (q *Queries) AddCaptcha(ctx context.Context, arg AddCaptchaParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, addCaptcha, arg.Email, arg.Captcha, arg.ExpiredAt)
}

const deleteCaptcha = `-- name: DeleteCaptcha :exec
UPDATE captcha
	SET is_deleted = 1
	WHERE id = ?
`

func (q *Queries) DeleteCaptcha(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteCaptcha, id)
	return err
}

const getCaptchaByEmailAndCode = `-- name: GetCaptchaByEmailAndCode :one
SELECT id, email, captcha, is_deleted, expired_at FROM captcha
	WHERE captcha = ?
	AND email = ?
	AND is_deleted = 0
	LIMIT 1
`

type GetCaptchaByEmailAndCodeParams struct {
	Captcha string `json:"captcha"`
	Email   string `json:"email"`
}

func (q *Queries) GetCaptchaByEmailAndCode(ctx context.Context, arg GetCaptchaByEmailAndCodeParams) (Captcha, error) {
	row := q.db.QueryRowContext(ctx, getCaptchaByEmailAndCode, arg.Captcha, arg.Email)
	var i Captcha
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Captcha,
		&i.IsDeleted,
		&i.ExpiredAt,
	)
	return i, err
}

const getCaptchaByID = `-- name: GetCaptchaByID :one
SELECT id, email, captcha, is_deleted, expired_at FROM captcha
	WHERE id = ?
	AND is_deleted = 0
	LIMIT 1
`

func (q *Queries) GetCaptchaByID(ctx context.Context, id int64) (Captcha, error) {
	row := q.db.QueryRowContext(ctx, getCaptchaByID, id)
	var i Captcha
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Captcha,
		&i.IsDeleted,
		&i.ExpiredAt,
	)
	return i, err
}

const todayEmailCount = `-- name: TodayEmailCount :one
SELECT COUNT(id) FROM captcha
WHERE email = ?
AND TO_DAYS(expired_at) = TO_DAYS(NOW())
`

func (q *Queries) TodayEmailCount(ctx context.Context, email string) (int64, error) {
	row := q.db.QueryRowContext(ctx, todayEmailCount, email)
	var count int64
	err := row.Scan(&count)
	return count, err
}

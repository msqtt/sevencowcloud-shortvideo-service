-- name: AddCaptcha :execresult
INSERT INTO captcha (
	email, captcha, expired_at
) VALUES ( ?, ?, ? );

-- name: TodayEmailCount :one
SELECT COUNT(id) FROM captcha
WHERE email = ?
AND TO_DAYS(expired_at) = TO_DAYS(NOW());

-- name: GetCaptchaByID :one
SELECT * FROM captcha
	WHERE id = ?
	AND is_deleted = 0
	LIMIT 1;

-- name: GetCaptchaByEmailAndCode :one
SELECT * FROM captcha
	WHERE captcha = ?
	AND email = ?
	AND is_deleted = 0
	LIMIT 1;

-- name: DeleteCaptcha :exec
UPDATE captcha
	SET is_deleted = 1
	WHERE id = ?;

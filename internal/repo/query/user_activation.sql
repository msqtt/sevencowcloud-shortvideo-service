-- name: AddActivation :execresult
INSERT INTO user_activation (
	user_id, activate_code, expired_at
) VALUES ( ?, ?, ? );

-- name: TodayActiationCount :one
SELECT COUNT(id) FROM user_activation 
WHERE user_id = ?
AND TO_DAYS(expired_at) = TO_DAYS(NOW());

-- name: GetActivationByID :one
SELECT * FROM user_activation
	WHERE id = ?
	AND is_deleted = 0
	LIMIT 1;

-- name: GetActivationByUserIDAndCode :one
SELECT * FROM user_activation
	WHERE user_id = ?
	AND activate_code = ?
	AND is_deleted = 0
	LIMIT 1;

-- name: ActivateUser :exec
UPDATE users
	SET is_activated = 1
	WHERE id = ?;

-- name: DeleteActivation :exec
UPDATE user_activation
	SET is_deleted = 1
	WHERE id = ?;

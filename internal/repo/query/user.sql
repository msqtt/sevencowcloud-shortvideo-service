-- name: AddUser :execresult
INSERT INTO users (
	nickname, email, password, profile_id
) VALUES ( 
	?, ?, ?, ?
);

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ? 
AND is_deleted = 0
LIMIT 1;

-- name: GetUserByNickName :one
SELECT * FROM users
WHERE nickname = ? 
AND is_deleted = 0
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ? 
AND is_deleted = 0
ORDER By created_at DESC
LIMIT 1;

-- name: UpdateNickName :execresult
UPDATE users
	SET nickname = ?
	WHERE id =?;

-- name: DeleteUser :exec
UPDATE users
	SET is_deleted = 1 
	WHERE id = ?;

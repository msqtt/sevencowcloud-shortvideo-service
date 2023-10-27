-- name: AddUser :execresult
INSERT INTO users (
	nickname, email, password, profile_id, is_deleted
) VALUES ( 
	?, ?, ?, ?, 0
);

-- name: GetUser :one
SELECT * FROM users
WHERE id = ? AND is_deleted = 0 LIMIT 1;

-- name: UpdateNickName :execresult
UPDATE users
	SET nickname = ?
	WHERE id =?;

-- name: DeleteUser :exec
UPDATE users
	SET is_deleted = 1 
	WHERE id = ?;

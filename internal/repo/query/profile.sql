-- name: AddProfile :execresult
INSERT INTO profiles (
	real_name, mood, gender, birth_date, introduction, avatar_link
) VALUES ( ?, ?, ?, ?, ?, ?);

-- name: GetProfile :one
SELECT * FROM profiles
WHERE id = ?;

-- name: UpdateProfile :exec
UPDATE profiles
	SET real_name = ?,
			mood = ?,
			gender = ?,
			birth_date = ?,
			introduction = ?,
			updated_at = now()
	WHERE id = ?;

-- name: UpdateAvatar :exec
UPDATE profiles
	SET avatar_link = ?
	WHERE id = ?;

-- name: DeleteProfile :exec
DELETE FROM profiles
	WHERE id = ?;

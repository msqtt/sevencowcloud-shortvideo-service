-- name: AddProfile :execresult
INSERT INTO profiles (
	real_name, mood, gender, birth_date, introduction
) VALUES ( ?, ?, ?, ?, ? );

-- name: GetProfile :one
SELECT * FROM profiles
WHERE id = ?;

-- name: UpdateProfile :execresult
UPDATE profiles
	SET real_name = ?,
			mood = ?,
			gender = ?,
			birth_date = ?,
			introduction = ?,
			updated_at = now()
	WHERE condition;

-- name: DeleteProfile :exec
DELETE FROM profiles
	WHERE id = ?;

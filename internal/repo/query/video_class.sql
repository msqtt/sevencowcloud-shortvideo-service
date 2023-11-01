-- name: GetVideoClassByID :one
SELECT * FROM video_class
WHERE id = ?
AND is_enabled = 1;

-- name: GetAllVideoClass :many
SELECT * FROM video_class
WHERE is_enabled = 1;

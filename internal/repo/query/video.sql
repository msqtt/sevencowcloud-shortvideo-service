-- name: AddVideo :execresult
INSERT INTO videos (
	content_hash, updated_at, cover_link, src_link	
) VALUES ( ?, now(), ?, ? );

-- name: GetVideoByID :one
SELECT * FROM videos
WHERE id = ?
LIMIT 1;

-- name: UpdateVideoLink :exec
UPDATE videos
	SET cover_link = ?, src_link = ?, updated_at = now()
	WHERE id = ?;

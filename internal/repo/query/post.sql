-- name: AddPost :execresult
INSERT INTO posts (
	title, description, user_id, video_id, created_at
) VALUES ( 
	?, ?, ?, ?, now()
);

-- name: TestGetAll :many
SELECT *, (SELECT count(id) FROM posts) total_size
FROM posts
WHERE is_deleted = 0
LIMIT ?, ?;

-- name: TestGetAllByTagID :many
SELECT p.*, (select count(*) from post_tag ipt where ipt.tag_id = ?) total_size  FROM post_tag pt
join posts p on pt.post_id = p.id
WHERE pt.tag_id = ?
AND is_deleted = 0
LIMIT ?, ?;

-- name: GetPostByID :one
SELECT * FROM posts
WHERE id = ?
AND is_deleted = 0
LIMIT 1;

-- name: GetPostByUserID :many
SELECT * FROM posts
WHERE user_id = ?
AND is_deleted = 0
LIMIT ?,?;

-- name: SearchPostByTitle :many
SELECT * FROM posts
WHERE title like "%?%"
AND is_deleted = 0
LIMIT ?,?;

-- name: UpdatePostInfo :exec
UPDATE posts
	SET title = ?, description = ?, updated_at = now()
	WHERE id = ?;

-- name: DeletePost :exec
UPDATE posts
	SET is_deleted = 1
	WHERE id = ?;

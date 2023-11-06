-- name: AddPostTag :exec
INSERT INTO post_tag (
	post_id, tag_id	
) VALUES ( 
	?, ?
);

-- name: GetTagByID :one
SELECT * FROM tags
WHERE id = ?
AND is_enabled = 1;

-- name: GetAllTags :many
SELECT * FROM tags
WHERE is_enabled = 1;

-- name: GetTagsByPostID :many
SELECT t.id, t.name, t.description  FROM tags t
join post_tag pt on pt.tag_id  = t.id
WHERE pt.post_id = ?;


-- name: SearchPostByTag :many
SELECT id FROM tags
WHERE tag_content like concat("%", ?,"%");

-- name: DeleteTag :exec
DELETE FROM post_tag
	WHERE post_id = ?
AND tag_id = ?;


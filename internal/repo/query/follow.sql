-- name: AddFollow :execresult
INSERT INTO follows (
	following_user_id, followed_user_id, followed_at
) VALUES ( ?, ?, now() );

-- name: GetFollow :one
SELECT * FROM follows
WHERE following_user_id = ?
AND followed_user_id = ?
LIMIT 1;

-- name: GetFollowingList :many
SELECT followed_at, u.id, u.nickname, u.email, u.created_at, p.real_name, p.mood, p.gender,
	p.birth_date, p.introduction, p.avatar_link
from follows f 
join users u on id = followed_user_id 
join profiles p on u.profile_id = p.id 
WHERE following_user_id = ?;

-- name: GetFollowedList :many
SELECT followed_at, u.id, u.nickname, u.email, u.created_at, p.real_name, p.mood, p.gender,
	p.birth_date, p.introduction, p.avatar_link
from follows f 
join users u on id = following_user_id 
join profiles p on u.profile_id = p.id 
WHERE followed_user_id = ?;

-- name: DeleteFollow :exec
DELETE FROM follows
	WHERE following_user_id = ?
	AND followed_user_id = ?;

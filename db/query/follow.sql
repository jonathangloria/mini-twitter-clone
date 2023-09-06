-- name: CreateFollowing :one
INSERT INTO follows (
  user_id,
  follower_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: ListFollower :many
SELECT * FROM follows
WHERE user_id = $1 LIMIT 20;

-- name: ListFollowing :many
SELECT * FROM follows
WHERE follower_id = $1 LIMIT 20;
-- name: CreateFollowing :one
INSERT INTO follows (
  user_id,
  follower_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: ListFollower :many
SELECT follower_id as id, users.username FROM follows 
INNER JOIN users ON users.id = follows.follower_id
WHERE user_id = $1 LIMIT 20;

-- name: ListFollowing :many
SELECT user_id as id, users.username FROM follows
INNER JOIN users ON users.id = follows.user_id
WHERE follower_id = $1 LIMIT 20;
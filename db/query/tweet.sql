-- name: CreateTweet :one
INSERT INTO tweets(
    user_id,
    body
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetTweet :one
SELECT tweets.id as tweet_id, users.id as user_id, users.username, tweets.body, tweets.created_at, tweets.edited_at 
FROM tweets INNER JOIN users
ON users.id = tweets.user_id
WHERE tweets.id = $1
LIMIT 1;

-- name: ListTweet :many
SELECT tweets.id as tweet_id, users.id as user_id, users.username, tweets.body, tweets.created_at, tweets.edited_at  
FROM tweets INNER JOIN users
ON users.id = tweets.user_id
WHERE tweets.user_id = $1 
LIMIT 10 OFFSET $2;

-- name: GetFeed :many
SELECT tweets.id as tweet_id, users.id as user_id, users.username, tweets.body, tweets.created_at, tweets.edited_at  
FROM tweets 
INNER JOIN users ON users.id = tweets.user_id
INNER JOIN follows ON users.id = follows.user_id
WHERE follows.follower_id = $1 ORDER BY tweets.created_at DESC 
LIMIT 10 OFFSET $2;

-- name: DeleteTweet :exec
DELETE FROM tweets
WHERE id = $1;

-- name: UpdateTweet :one
UPDATE tweets
SET
  body = $1,
  edited_at = $2 
WHERE
  id = $3
RETURNING *;
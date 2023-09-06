-- name: CreateUser :one
INSERT INTO users (
  email,
  username,
  passhash,
  full_name
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
  passhash = COALESCE(sqlc.narg(passhash), passhash),
  full_name = COALESCE(sqlc.narg(full_name), full_name),
  username = COALESCE(sqlc.narg(username), username)
WHERE
  id = sqlc.arg(id)
RETURNING *;
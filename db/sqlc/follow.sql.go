// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: follow.sql

package db

import (
	"context"
)

const createFollowing = `-- name: CreateFollowing :one
INSERT INTO follows (
  user_id,
  follower_id
) VALUES (
  $1, $2
) RETURNING user_id, follower_id
`

type CreateFollowingParams struct {
	UserID     int64 `json:"user_id"`
	FollowerID int64 `json:"follower_id"`
}

func (q *Queries) CreateFollowing(ctx context.Context, arg CreateFollowingParams) (Follow, error) {
	row := q.db.QueryRowContext(ctx, createFollowing, arg.UserID, arg.FollowerID)
	var i Follow
	err := row.Scan(&i.UserID, &i.FollowerID)
	return i, err
}

const listFollower = `-- name: ListFollower :many
SELECT user_id, follower_id FROM follows
WHERE user_id = $1 LIMIT 20
`

func (q *Queries) ListFollower(ctx context.Context, userID int64) ([]Follow, error) {
	rows, err := q.db.QueryContext(ctx, listFollower, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Follow{}
	for rows.Next() {
		var i Follow
		if err := rows.Scan(&i.UserID, &i.FollowerID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listFollowing = `-- name: ListFollowing :many
SELECT user_id, follower_id FROM follows
WHERE follower_id = $1 LIMIT 20
`

func (q *Queries) ListFollowing(ctx context.Context, followerID int64) ([]Follow, error) {
	rows, err := q.db.QueryContext(ctx, listFollowing, followerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Follow{}
	for rows.Next() {
		var i Follow
		if err := rows.Scan(&i.UserID, &i.FollowerID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package db

import (
	"context"
)

type Querier interface {
	CreateFollowing(ctx context.Context, arg CreateFollowingParams) (Follow, error)
	CreateTweet(ctx context.Context, arg CreateTweetParams) (Tweet, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteTweet(ctx context.Context, id int64) error
	GetFeed(ctx context.Context, arg GetFeedParams) ([]GetFeedRow, error)
	GetTweet(ctx context.Context, id int64) (GetTweetRow, error)
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	ListFollower(ctx context.Context, userID int64) ([]ListFollowerRow, error)
	ListFollowing(ctx context.Context, followerID int64) ([]ListFollowingRow, error)
	ListTweet(ctx context.Context, arg ListTweetParams) ([]ListTweetRow, error)
	UpdateTweet(ctx context.Context, arg UpdateTweetParams) (Tweet, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)

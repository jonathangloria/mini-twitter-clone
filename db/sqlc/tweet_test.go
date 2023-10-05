package db

import (
	"context"
	"testing"
	"time"

	"github.com/jonathangloria/mini-twitter-clone/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomTweet(t *testing.T, userID int64) Tweet {
	body := util.RandomString(30)
	arg := CreateTweetParams{
		UserID: userID,
		Body:   body,
	}
	tweet, err := testQueries.CreateTweet(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tweet)

	require.Equal(t, userID, tweet.UserID)
	require.Equal(t, body, tweet.Body)
	require.NotZero(t, tweet.ID)
	require.NotZero(t, tweet.CreatedAt)
	require.Zero(t, tweet.EditedAt)

	return tweet
}

func TestCreateTweet(t *testing.T) {
	user := createRandomUser(t)
	CreateRandomTweet(t, user.ID)
}

func TestGetTweet(t *testing.T) {
	user := createRandomUser(t)
	exptweet := CreateRandomTweet(t, user.ID)
	tweet, err := testQueries.GetTweet(context.Background(), exptweet.ID)
	require.NoError(t, err)
	require.NotEmpty(t, tweet)

	require.Equal(t, exptweet.ID, tweet.TweetID)
	require.Equal(t, exptweet.Body, tweet.Body)
	require.WithinDuration(t, exptweet.CreatedAt, tweet.CreatedAt, time.Second)
	require.WithinDuration(t, exptweet.EditedAt, tweet.EditedAt, time.Second)
}

func TestListTweet(t *testing.T) {
	user1 := createRandomUser(t)
	for i := 0; i < 5; i++ {
		CreateRandomTweet(t, user1.ID)
	}

	tweets, err := testQueries.ListTweet(context.Background(), ListTweetParams{
		UserID: user1.ID,
		Offset: 0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, tweets)

	require.Len(t, tweets, 5)

	for _, tweet := range tweets {
		require.NotEmpty(t, tweet)
		require.Equal(t, user1.ID, tweet.UserID)
	}
}

func TestGetFeed(t *testing.T) {
	user, following := CreateListFollowing(t)
	listFollowing := make([]string, 5)

	for _, user2 := range following {
		CreateRandomTweet(t, user2.FollowingID)
		listFollowing = append(listFollowing, user2.Username)
	}

	feed, err := testQueries.GetFeed(context.Background(), GetFeedParams{
		FollowerID: user.ID,
		Offset:     0,
	})

	require.NoError(t, err)
	require.NotEmpty(t, feed)

	require.Len(t, feed, 5)

	for _, tweet := range feed {
		require.NotEmpty(t, tweet)
		require.Contains(t, listFollowing, tweet.Username)
	}
}

func TestDeleteTweet(t *testing.T) {
	user := createRandomUser(t)
	tweet := CreateRandomTweet(t, user.ID)

	err := testQueries.DeleteTweet(context.Background(), tweet.ID)
	require.NoError(t, err)
	_, err = testQueries.GetTweet(context.Background(), tweet.ID)
	require.Error(t, err)
}

func TestUpdateTweet(t *testing.T) {
	user := createRandomUser(t)
	oldtweet := CreateRandomTweet(t, user.ID)
	newbody := util.RandomString(25)

	newtweet, err := testQueries.UpdateTweet(context.Background(), UpdateTweetParams{
		Body:     newbody,
		EditedAt: time.Now(),
		ID:       oldtweet.ID,
	})

	require.NoError(t, err)
	require.NotEmpty(t, newtweet)

	require.Equal(t, oldtweet.ID, newtweet.ID)
	require.Equal(t, newtweet.Body, newbody)
	require.NotZero(t, newtweet.EditedAt)
	require.WithinDuration(t, newtweet.EditedAt, time.Now(), time.Second)
	require.WithinDuration(t, oldtweet.CreatedAt, newtweet.CreatedAt, time.Second)
}

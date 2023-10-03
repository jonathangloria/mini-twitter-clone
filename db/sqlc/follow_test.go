package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomFollower(t *testing.T, expectedUser User) {
	expectedFollower := createRandomUser(t)
	follow, err := testQueries.CreateFollowing(context.Background(), CreateFollowingParams{
		UserID:     expectedUser.ID,
		FollowerID: expectedFollower.ID,
	})

	require.NoError(t, err)
	require.NotEmpty(t, expectedUser)
	require.NotEmpty(t, expectedFollower)
	require.NotEmpty(t, follow)

	require.Equal(t, expectedUser.ID, follow.UserID)
	require.Equal(t, expectedFollower.ID, follow.FollowerID)
}

func createRandomFollowing(t *testing.T, expectedFollower User) {
	expectedUser := createRandomUser(t)
	follow, err := testQueries.CreateFollowing(context.Background(), CreateFollowingParams{
		UserID:     expectedUser.ID,
		FollowerID: expectedFollower.ID,
	})

	require.NoError(t, err)
	require.NotEmpty(t, expectedUser)
	require.NotEmpty(t, expectedFollower)
	require.NotEmpty(t, follow)

	require.Equal(t, expectedUser.ID, follow.UserID)
	require.Equal(t, expectedFollower.ID, follow.FollowerID)
}

func TestCreateFollowing(t *testing.T) {
	user1 := createRandomUser(t)
	createRandomFollower(t, user1)
}

func TestListFollower(t *testing.T) {
	user1 := createRandomUser(t)
	for i := 0; i < 5; i++ {
		createRandomFollower(t, user1)
	}

	followers, err := testQueries.ListFollower(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, followers)

	require.Len(t, followers, 5)

	for _, follower := range followers {
		require.NotEmpty(t, follower)
		require.Equal(t, user1.ID, follower.UserID)
	}
}

func TestListFollowing(t *testing.T) {
	user1 := createRandomUser(t)
	for i := 0; i < 5; i++ {
		createRandomFollowing(t, user1)
	}

	following, err := testQueries.ListFollowing(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, following)

	require.Len(t, following, 5)

	for _, user := range following {
		require.NotEmpty(t, user)
		require.Equal(t, user1.ID, user.UserID)
	}
}

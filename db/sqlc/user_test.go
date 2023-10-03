package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jonathangloria/mini-twitter-clone/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	passhass, _ := util.HashPassword(util.RandomString(6))
	arg := CreateUserParams{
		Email:    util.RandomEmail(),
		Username: util.RandomUsername(),
		Passhash: passhass,
		FullName: util.RandomFullname(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Passhash, user.Passhash)
	require.Equal(t, arg.FullName, user.FullName)

	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	expected_user := createRandomUser(t)
	actual_user, err := testQueries.GetUser(context.Background(), expected_user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, actual_user)

	require.Equal(t, expected_user.Username, actual_user.Username)
	require.Equal(t, expected_user.Passhash, actual_user.Passhash)
	require.Equal(t, expected_user.FullName, actual_user.FullName)
	require.Equal(t, expected_user.Email, actual_user.Email)
	require.WithinDuration(t, expected_user.CreatedAt, actual_user.CreatedAt, time.Second)
}

func TestGetUserByUsername(t *testing.T) {
	expected_user := createRandomUser(t)
	actual_user, err := testQueries.GetUserByUsername(context.Background(), expected_user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, actual_user)

	require.Equal(t, expected_user.Username, actual_user.Username)
	require.Equal(t, expected_user.Passhash, actual_user.Passhash)
	require.Equal(t, expected_user.FullName, actual_user.FullName)
	require.Equal(t, expected_user.Email, actual_user.Email)
	require.WithinDuration(t, expected_user.CreatedAt, actual_user.CreatedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullname := util.RandomFullname()

	arg := UpdateUserParams{
		ID: oldUser.ID,
		FullName: sql.NullString{
			String: newFullname,
			Valid:  true,
		},
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newFullname, updatedUser.FullName)
	require.Equal(t, oldUser.Passhash, updatedUser.Passhash)
	require.Equal(t, oldUser.Email, updatedUser.Email)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)
	newEmail := util.RandomEmail()

	arg := UpdateUserParams{
		ID: oldUser.ID,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, oldUser.Passhash, updatedUser.Passhash)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)
	newpass := util.RandomString(6)
	newPasshash, err := util.HashPassword(newpass)
	require.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		ID: oldUser.ID,
		Passhash: sql.NullString{
			String: newPasshash,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Passhash, updatedUser.Passhash)
	require.Equal(t, newPasshash, updatedUser.Passhash)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
}

func TestUdpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)
	newpass := util.RandomString(6)
	newPasshash, err := util.HashPassword(newpass)
	require.NoError(t, err)
	newEmail := util.RandomEmail()
	newFullname := util.RandomFullname()

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		ID: oldUser.ID,
		Passhash: sql.NullString{
			String: newPasshash,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullname,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Passhash, updatedUser.Passhash)
	require.Equal(t, newPasshash, updatedUser.Passhash)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newFullname, updatedUser.FullName)
}

package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vladoohr/simple_bank/util"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	createUserParams := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), createUserParams)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, createUserParams.Username, user.Username)
	require.Equal(t, createUserParams.HashedPassword, user.HashedPassword)
	require.Equal(t, createUserParams.FullName, user.FullName)
	require.Equal(t, createUserParams.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.PasswordChangeAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, 1*time.Second)
	require.WithinDuration(t, user1.PasswordChangeAt, user2.PasswordChangeAt, 1*time.Second)
}

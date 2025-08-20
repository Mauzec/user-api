package db

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/mauzec/user-api/internal/util"
	"github.com/stretchr/testify/assert"
)

func createAndTestRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(
		util.RandomString(10),
	)

	assert.NoError(t, err)
	args := CreateUserParams{
		Username:       util.RandomUsername(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
		FullName:       fmt.Sprintf("%s %s", util.RandomString(5), util.RandomString(5)),
		Phone:          util.RandomPhone(),
		Gender:         "M",
		Age:            int32(util.RandomInt(18, 60)),
	}

	ctx := context.Background()
	user, err := testQueries.CreateUser(ctx, args)
	if !assert.NoError(t, err) {
		log.Fatal(
			fmt.Errorf(
				"unable to create user with args:\n %+v\n%v", args, err),
		)
	}

	if !assert.NotEmpty(t, user) {
		log.Fatal(
			fmt.Errorf("received empty user, but args:\n%+v", args),
		)
	}

	assert.Equal(t, args.Username, user.Username)
	assert.Equal(t, args.Email, user.Email)
	assert.Equal(t, args.FullName, user.FullName)
	assert.Equal(t, args.HashedPassword, user.HashedPassword)
	assert.Equal(t, args.Phone, user.Phone)

	assert.True(t, user.PasswordChangedAt.Time.IsZero())
	assert.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	_ = createAndTestRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createAndTestRandomUser(t)

	gotUser, err := testQueries.GetUserByID(context.Background(), user.ID)
	assert.NoError(t, err)

	assert.Equal(t, user.ID, gotUser.ID)
	assert.Equal(t, user.Username, gotUser.Username)
	assert.Equal(t, user.Email, gotUser.Email)
	assert.Equal(t, user.FullName, gotUser.FullName)
	assert.Equal(t, user.HashedPassword, gotUser.HashedPassword)
	assert.Equal(t, user.Phone, gotUser.Phone)
	assert.Equal(t, user.PasswordChangedAt, gotUser.PasswordChangedAt)
	assert.Equal(t, user.CreatedAt, gotUser.CreatedAt)
}

func TestGetUserByUsername(t *testing.T) {
	user := createAndTestRandomUser(t)

	gotUser, err := testQueries.GetUserByUsername(context.Background(), user.Username)
	assert.NoError(t, err)

	assert.Equal(t, user.ID, gotUser.ID)
	assert.Equal(t, user.Username, gotUser.Username)
	assert.Equal(t, user.Email, gotUser.Email)
	assert.Equal(t, user.FullName, gotUser.FullName)
	assert.Equal(t, user.HashedPassword, gotUser.HashedPassword)
	assert.Equal(t, user.Phone, gotUser.Phone)
	assert.Equal(t, user.PasswordChangedAt, gotUser.PasswordChangedAt)
	assert.Equal(t, user.CreatedAt, gotUser.CreatedAt)
}

func TestUpdateUser(t *testing.T) {
	var user User

	t.Run("Creating new user", func(t *testing.T) {
		user = createAndTestRandomUser(t)
	})

	t.Run("updating username", func(t *testing.T) {
		wantArgs := UpdateUserParams{
			ID:       user.ID,
			Email:    "t@t.com",
			FullName: "Ka Ma",
			Gender:   "M",
			Phone:    "+123413424",
		}

		newUser, err := testQueries.UpdateUser(context.Background(), wantArgs)
		if !assert.NoError(t, err) {
			log.Fatal(fmt.Errorf("error updating user with args:\n%+v\n%v", wantArgs, err))
		}
		if !assert.NotEmpty(t, newUser) {
			log.Fatal(fmt.Errorf("received empty user, but args:\n%+v", wantArgs))
		}

		assert.Equal(t, wantArgs.ID, newUser.ID)
		assert.Equal(t, user.Username, newUser.Username)
		assert.Equal(t, wantArgs.Email, newUser.Email)
		assert.Equal(t, wantArgs.FullName, newUser.FullName)
		assert.Equal(t, user.HashedPassword, newUser.HashedPassword)
		assert.Equal(t, wantArgs.Phone, newUser.Phone)
		assert.Equal(t, user.PasswordChangedAt, newUser.PasswordChangedAt)
		assert.Equal(t, user.CreatedAt, newUser.CreatedAt)

		user = newUser
	})

}

func TestDeleteUser(t *testing.T) {
	var user User

	t.Run("creating new user", func(t *testing.T) {
		user = createAndTestRandomUser(t)
	})

	t.Run("deleting user", func(t *testing.T) {
		err := testQueries.DeleteUserByID(context.Background(), user.ID)
		if !assert.NoError(t, err) {
			log.Fatal(fmt.Errorf("error deleting user with ID: %d, error: %v", user.ID, err))
		}
	})

	t.Run("getting deleted user", func(t *testing.T) {
		newUser, err := testQueries.GetUserByID(context.Background(), user.ID)
		assert.Error(t, err)
		assert.EqualError(t, err, pgx.ErrNoRows.Error())
		assert.Empty(t, newUser)
	})
}

package service_test

import (
	"context"
	"github.com/OzkrOssa/freeradius-api/internal/core/domain"
	"github.com/OzkrOssa/freeradius-api/internal/core/port/mocks"
	"github.com/OzkrOssa/freeradius-api/internal/core/service"
	"github.com/OzkrOssa/freeradius-api/internal/core/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type registerInput struct {
	user *domain.User
}

type expectedOutput struct {
	user *domain.User
	err  error
}

func TestUserService_CreateUser(t *testing.T) {
	ctx := context.Background()
	email := gofakeit.Email()
	name := gofakeit.Name()
	password := gofakeit.Password(true, true, true, true, true, 10)
	hashedPassword, _ := utils.HashPassword(password)

	userInput := &domain.User{
		Email:    email,
		Name:     name,
		Password: password,
	}
	userOutput := &domain.User{
		ID:        gofakeit.Uint64(),
		Email:     email,
		Name:      name,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	serializedUser, _ := utils.Serialize(userOutput)
	cacheKey := utils.GenerateCacheKey("user", userOutput.ID)
	ttl := time.Duration(0)

	testCases := []struct {
		desc     string
		mocks    func(repo *mocks.UserRepository, cache *mocks.CacheRepository)
		input    registerInput
		expected expectedOutput
	}{
		{
			desc: "Success",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("CreateUser", ctx, userInput).Return(userOutput, nil)
				cache.On("Set", ctx, cacheKey, serializedUser, ttl).Return(nil)
				cache.On("DeleteByPrefix", ctx, "users:*").Return(nil)
			},
			input: registerInput{user: userInput},
			expected: expectedOutput{
				user: userOutput,
				err:  nil,
			},
		},
		{
			desc: "Fail_DuplicateData",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("CreateUser", ctx, userInput).Return(nil, domain.ConflictDataError)
			},
			input: registerInput{user: userInput},
			expected: expectedOutput{
				user: nil,
				err:  domain.ConflictDataError,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("CreateUser", ctx, userInput).Return(nil, domain.InternalError)
			},
			input: registerInput{user: userInput},
			expected: expectedOutput{
				user: nil,
				err:  domain.InternalError,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("CreateUser", ctx, userInput).Return(userOutput, nil)
				cache.On("Set", ctx, cacheKey, serializedUser, ttl).Return(domain.InternalError)
			},
			input: registerInput{user: userInput},
			expected: expectedOutput{
				user: nil,
				err:  domain.InternalError,
			},
		},
		{
			desc: "Fail_DeleteCacheByPrefix",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("CreateUser", ctx, userInput).Return(userOutput, nil)
				cache.On("Set", ctx, cacheKey, serializedUser, ttl).Return(nil)
				cache.On("DeleteByPrefix", ctx, "users:*").Return(domain.InternalError)
			},
			input: registerInput{user: userInput},
			expected: expectedOutput{
				user: nil,
				err:  domain.InternalError,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)
			cache := mocks.NewCacheRepository(t)
			tc.mocks(repo, cache)

			userService := service.NewUserService(repo, cache)

			user, err := userService.CreateUser(ctx, tc.input.user)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.user, user, "User mismatch")
		})
	}

}

type getUserTestedInput struct {
	ID uint64
}

type getUserExpectedOutput struct {
	user *domain.User
	err  error
}

func TestUserService_GetUser(t *testing.T) {
	ctx := context.Background()
	id := gofakeit.Uint64()

	userOutput := &domain.User{
		ID:       id,
		Email:    gofakeit.Email(),
		Name:     gofakeit.Name(),
		Password: gofakeit.Password(true, true, true, true, true, 10),
	}

	cacheKey := utils.GenerateCacheKey("user", id)
	userSerialized, _ := utils.Serialize(userOutput)
	ttl := time.Duration(0)

	testCases := []struct {
		desc     string
		mocks    func(repo *mocks.UserRepository, cache *mocks.CacheRepository)
		input    getUserTestedInput
		expected getUserExpectedOutput
	}{
		{
			desc: "Success_FromCache",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				cache.On("Get", ctx, cacheKey).Return(userSerialized, nil)
			},
			input: getUserTestedInput{ID: id},
			expected: getUserExpectedOutput{
				user: userOutput,
				err:  nil,
			},
		},
		{
			desc: "Success_FromDB",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				cache.On("Get", ctx, cacheKey).Return(nil, domain.DataNotFoundError)
				repo.On("GetUserById", ctx, id).Return(userOutput, nil)
				cache.On("Set", ctx, cacheKey, userSerialized, ttl).Return(nil)
			},
			input: getUserTestedInput{ID: id},
			expected: getUserExpectedOutput{
				user: userOutput,
				err:  nil,
			},
		},
		{
			desc: "Fail_NotFound",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				cache.On("Get", ctx, cacheKey).Return(nil, domain.DataNotFoundError)
				repo.On("GetUserById", ctx, id).Return(nil, domain.DataNotFoundError)
			},
			input: getUserTestedInput{ID: id},
			expected: getUserExpectedOutput{
				user: nil,
				err:  domain.DataNotFoundError,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				cache.On("Get", ctx, cacheKey).Return(nil, domain.InternalError)
				repo.On("GetUserById", ctx, id).Return(nil, domain.InternalError)
			},
			input: getUserTestedInput{ID: id},
			expected: getUserExpectedOutput{
				user: nil,
				err:  domain.InternalError,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				cache.On("Get", ctx, cacheKey).Return(nil, domain.DataNotFoundError)
				repo.On("GetUserById", ctx, id).Return(userOutput, nil)
				cache.On("Set", ctx, cacheKey, userSerialized, ttl).Return(domain.InternalError)
			},
			input: getUserTestedInput{ID: id},
			expected: getUserExpectedOutput{
				user: nil,
				err:  domain.InternalError,
			},
		},
		{
			desc: "Fail_Deserialize",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				cache.On("Get", ctx, cacheKey).Return([]byte("bat user"), nil)
			},
			input: getUserTestedInput{ID: id},
			expected: getUserExpectedOutput{
				user: nil,
				err:  domain.InternalError,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)
			cache := mocks.NewCacheRepository(t)
			tc.mocks(repo, cache)

			userService := service.NewUserService(repo, cache)

			user, err := userService.GetUser(ctx, id)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.user, user, "User mismatch")
		})
	}
}

type listUsersTestedInput struct {
	skip  uint64
	limit uint64
}

type listUsersExpectedOutput struct {
	users []domain.User
	err   error
}

func TestUserService_ListUsers(t *testing.T) {
	ctx := context.Background()
	skip, limit := gofakeit.Uint64(), gofakeit.Uint64()
	params := utils.GenerateCacheKeyParams(skip, limit)
	cacheKey := utils.GenerateCacheKey("users", params)

	var users []domain.User

	for i := 0; i < 10; i++ {
		userPassword := gofakeit.Password(true, true, true, true, false, 8)
		hashedPassword, _ := utils.HashPassword(userPassword)

		users = append(users, domain.User{
			ID:       gofakeit.Uint64(),
			Name:     gofakeit.Name(),
			Email:    gofakeit.Email(),
			Password: hashedPassword,
		})
	}

	usersSerialized, _ := utils.Serialize(users)
	ttl := time.Duration(0)

	testCases := []struct {
		desc     string
		mocks    func(repo *mocks.UserRepository, cache *mocks.CacheRepository)
		input    listUsersTestedInput
		expected listUsersExpectedOutput
	}{
		{
			desc: "Success_FromCache",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				cache.On("Get", ctx, cacheKey).Return(usersSerialized, nil)
			},
			input: listUsersTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listUsersExpectedOutput{
				users: users,
				err:   nil,
			},
		},
		{
			desc: "Success_FromDB",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				cache.On("Get", ctx, cacheKey).Return(nil, domain.DataNotFoundError)
				repo.On("ListUsers", ctx, skip, limit).Return(users, nil)
				cache.On("Set", ctx, cacheKey, usersSerialized, ttl).Return(nil)
			},
			input: listUsersTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listUsersExpectedOutput{
				users: users,
				err:   nil,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				cache.On("Get", ctx, cacheKey).Return(nil, domain.InternalError)
				repo.On("ListUsers", ctx, skip, limit).Return(nil, domain.InternalError)
			},
			input: listUsersTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listUsersExpectedOutput{
				users: nil,
				err:   domain.InternalError,
			},
		},
		{
			desc: "Fail_Deserialize",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				cache.On("Get", ctx, cacheKey).Return([]byte("invalid"), nil)
			},
			input: listUsersTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listUsersExpectedOutput{
				users: nil,
				err:   domain.InternalError,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				cache.On("Get", ctx, cacheKey).Return(nil, domain.DataNotFoundError)
				repo.On("ListUsers", ctx, skip, limit).Return(users, nil)
				cache.On("Set", ctx, cacheKey, usersSerialized, ttl).Return(domain.InternalError)
			},
			input: listUsersTestedInput{
				skip:  skip,
				limit: limit,
			},
			expected: listUsersExpectedOutput{
				users: nil,
				err:   domain.InternalError,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)
			cache := mocks.NewCacheRepository(t)
			tc.mocks(repo, cache)

			userService := service.NewUserService(repo, cache)

			users, err := userService.ListUsers(ctx, tc.input.skip, tc.input.limit)
			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.users, users, "Users mismatch")
		})
	}
}

type updateUserTestedInput struct {
	user *domain.User
}

type updateUserExpectedOutput struct {
	user *domain.User
	err  error
}

func TestUserService_UpdateUser(t *testing.T) {
	ctx := context.Background()
	id := gofakeit.Uint64()

	userInput := &domain.User{
		ID:    id,
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
	}

	userOutput := &domain.User{
		ID:    id,
		Name:  userInput.Name,
		Email: userInput.Email,
	}

	existingUser := &domain.User{
		ID:    id,
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
	}

	cacheKey := utils.GenerateCacheKey("user", id)
	userSerialized, _ := utils.Serialize(userOutput)
	ttl := time.Duration(0)

	testCases := []struct {
		desc     string
		mocks    func(repo *mocks.UserRepository, cache *mocks.CacheRepository)
		input    updateUserTestedInput
		expected updateUserExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("GetUserById", ctx, id).Return(existingUser, nil)
				repo.On("UpdateUser", ctx, userInput).Return(userOutput, nil)
				cache.On("Delete", ctx, cacheKey).Return(nil)
				cache.On("Set", ctx, cacheKey, userSerialized, ttl).Return(nil)
				cache.On("DeleteByPrefix", ctx, "users:*").Return(nil)

			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: userOutput,
				err:  nil,
			},
		},
		{
			desc: "Fail_NotFound",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("GetUserById", ctx, id).Return(nil, domain.DataNotFoundError)

			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  domain.DataNotFoundError,
			},
		},
		{
			desc: "Fail_InternalErrorGetById",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("GetUserById", ctx, id).Return(nil, domain.InternalError)

			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  domain.InternalError,
			},
		},
		{
			desc: "Fail_EmptyData",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("GetUserById", ctx, id).Return(existingUser, nil)
			},
			input: updateUserTestedInput{
				user: &domain.User{
					ID: id,
				},
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  domain.NoUpdatedDataError,
			},
		},
		{
			desc: "Fail_SameData",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("GetUserById", ctx, id).Return(existingUser, nil)
			},
			input: updateUserTestedInput{
				user: existingUser,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  domain.NoUpdatedDataError,
			},
		},
		{
			desc: "Fail_DuplicateData",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("GetUserById", ctx, id).Return(existingUser, nil)
				repo.On("UpdateUser", ctx, userInput).Return(nil, domain.ConflictDataError)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  domain.ConflictDataError,
			},
		},
		{
			desc: "Fail_InternalErrorUpdate",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("GetUserById", ctx, id).Return(existingUser, nil)
				repo.On("UpdateUser", ctx, userInput).Return(nil, domain.InternalError)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  domain.InternalError,
			},
		},
		{
			desc: "Fail_DeleteCache",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("GetUserById", ctx, id).Return(existingUser, nil)
				repo.On("UpdateUser", ctx, userInput).Return(userOutput, nil)
				cache.On("Delete", ctx, cacheKey).Return(domain.InternalError)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  domain.InternalError,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("GetUserById", ctx, id).Return(existingUser, nil)
				repo.On("UpdateUser", ctx, userInput).Return(userOutput, nil)
				cache.On("Delete", ctx, cacheKey).Return(nil)
				cache.On("Set", ctx, cacheKey, userSerialized, ttl).Return(domain.InternalError)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  domain.InternalError,
			},
		},
		{
			desc: "Fail_DeleteByPrefix",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("GetUserById", ctx, id).Return(existingUser, nil)
				repo.On("UpdateUser", ctx, userInput).Return(userOutput, nil)
				cache.On("Delete", ctx, cacheKey).Return(nil)
				cache.On("Set", ctx, cacheKey, userSerialized, ttl).Return(nil)
				cache.On("DeleteByPrefix", ctx, "users:*").Return(domain.InternalError)
			},
			input: updateUserTestedInput{
				user: userInput,
			},
			expected: updateUserExpectedOutput{
				user: nil,
				err:  domain.InternalError,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)
			cache := mocks.NewCacheRepository(t)
			tc.mocks(repo, cache)
			userService := service.NewUserService(repo, cache)

			user, err := userService.UpdateUser(ctx, tc.input.user)

			assert.Equal(t, tc.expected.err, err, "Error mismatch")
			assert.Equal(t, tc.expected.user, user, "Users mismatch")
		})
	}
}

type userDeleteExpectedOutput struct {
	err error
}

func TestUserService_DeleteUser(t *testing.T) {
	ctx := context.Background()
	id := gofakeit.Uint64()

	cacheKey := utils.GenerateCacheKey("user", id)

	testCases := []struct {
		desc     string
		mocks    func(repo *mocks.UserRepository, cache *mocks.CacheRepository)
		input    uint64
		expected userDeleteExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("GetUserById", ctx, id).Return(&domain.User{}, nil)
				cache.On("Delete", ctx, cacheKey).Return(nil)
				cache.On("DeleteByPrefix", ctx, "users:*").Return(nil)
				repo.On("DeleteUser", ctx, id).Return(nil)
			},
			input: id,
			expected: userDeleteExpectedOutput{
				err: nil,
			},
		},
		{
			desc: "Fail_NotFound",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("GetUserById", ctx, id).Return(nil, domain.DataNotFoundError)
			},
			input: id,
			expected: userDeleteExpectedOutput{
				err: domain.DataNotFoundError,
			},
		},
		{
			desc: "Fail_InternalErrorGetByID",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("GetUserById", ctx, id).Return(nil, domain.InternalError)
			},
			input: id,
			expected: userDeleteExpectedOutput{
				err: domain.InternalError,
			},
		},
		{
			desc: "Fail_DeleteCache",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("GetUserById", ctx, id).Return(&domain.User{}, nil)
				cache.On("Delete", ctx, cacheKey).Return(domain.InternalError)
			},
			input: id,
			expected: userDeleteExpectedOutput{
				err: domain.InternalError,
			},
		},
		{
			desc: "Fail_DeleteByPrefix",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("GetUserById", ctx, id).Return(&domain.User{}, nil)
				cache.On("Delete", ctx, cacheKey).Return(nil)
				cache.On("DeleteByPrefix", ctx, "users:*").Return(domain.InternalError)
			},
			input: id,
			expected: userDeleteExpectedOutput{
				err: domain.InternalError,
			},
		},
		{
			desc: "Fail_InternalErrorDelete",
			mocks: func(repo *mocks.UserRepository, cache *mocks.CacheRepository) {
				repo.On("GetUserById", ctx, id).Return(&domain.User{}, nil)
				cache.On("Delete", ctx, cacheKey).Return(nil)
				cache.On("DeleteByPrefix", ctx, "users:*").Return(nil)
				repo.On("DeleteUser", ctx, id).Return(domain.InternalError)
			},
			input: id,
			expected: userDeleteExpectedOutput{
				err: domain.InternalError,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)
			cache := mocks.NewCacheRepository(t)
			tc.mocks(repo, cache)
			userService := service.NewUserService(repo, cache)

			err := userService.DeleteUser(ctx, tc.input)

			assert.Equal(t, tc.expected.err, err, "Error mismatch")
		})
	}

}

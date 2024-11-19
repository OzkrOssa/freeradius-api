package service

import (
	"context"
	"github.com/OzkrOssa/freeradius-api/internal/core/domain"
	"github.com/OzkrOssa/freeradius-api/internal/core/port/mocks"
	"github.com/OzkrOssa/freeradius-api/internal/core/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type LoginTestedInput struct {
	email    string
	password string
}

type LoginTestedOutput struct {
	token string
	err   error
}

func TestAuthService_Login(t *testing.T) {
	ctx := context.Background()

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, true, 8)

	simulateToken := gofakeit.UUID()
	hashPassword, _ := utils.HashPassword(password)

	userOutput := &domain.User{
		ID:        gofakeit.Uint64(),
		Email:     email,
		Password:  hashPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	failUser := &domain.User{
		ID:       gofakeit.Uint64(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 8),
	}

	testCases := []struct {
		desc     string
		mocks    func(repo *mocks.UserRepository, ts *mocks.TokenService)
		input    LoginTestedInput
		expected LoginTestedOutput
	}{
		{
			desc: "Success",
			mocks: func(repo *mocks.UserRepository, ts *mocks.TokenService) {
				repo.On("GetUserByEmail", ctx, email).Return(userOutput, nil)
				ts.On("CreateToken", userOutput).Return(simulateToken, nil)
			},
			input: LoginTestedInput{
				email:    email,
				password: password,
			},
			expected: LoginTestedOutput{
				token: simulateToken,
				err:   nil,
			},
		},
		{
			desc: "Fail_UserNotFound",
			mocks: func(repo *mocks.UserRepository, ts *mocks.TokenService) {
				repo.On("GetUserByEmail", ctx, email).Return(nil, domain.DataNotFoundError)
			},
			input: LoginTestedInput{
				email:    email,
				password: password,
			},
			expected: LoginTestedOutput{
				token: "",
				err:   domain.InvalidCredentialsError,
			},
		},
		{
			desc: "Fail_PasswordMismatch",
			mocks: func(repo *mocks.UserRepository, ts *mocks.TokenService) {
				repo.On("GetUserByEmail", ctx, email).Return(failUser, nil)
			},
			input: LoginTestedInput{
				email:    email,
				password: password,
			},
			expected: LoginTestedOutput{
				token: "",
				err:   domain.InvalidCredentialsError,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(repo *mocks.UserRepository, ts *mocks.TokenService) {
				repo.On("GetUserByEmail", ctx, email).Return(nil, domain.InternalError)
			},
			input: LoginTestedInput{
				email:    email,
				password: password,
			},
			expected: LoginTestedOutput{
				token: "",
				err:   domain.InternalError,
			},
		},
		{
			desc: "Fail_InternalErrorCreateToken",
			mocks: func(repo *mocks.UserRepository, ts *mocks.TokenService) {
				repo.On("GetUserByEmail", ctx, email).Return(userOutput, nil)
				ts.On("CreateToken", userOutput).Return("", domain.TokenCreationError)
			},
			input: LoginTestedInput{
				email:    email,
				password: password,
			},
			expected: LoginTestedOutput{
				token: "",
				err:   domain.TokenCreationError,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			repo := mocks.NewUserRepository(t)
			ts := mocks.NewTokenService(t)

			tc.mocks(repo, ts)

			authService := NewAuthService(repo, ts)

			token, err := authService.Login(ctx, tc.input.email, tc.input.password)

			assert.Equal(t, tc.expected.err, err, "Error Mismatch")
			assert.Equal(t, tc.expected.token, token, "Token Mismatch")
		})
	}
}

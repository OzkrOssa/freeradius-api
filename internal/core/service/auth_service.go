package service

import (
	"context"
	"errors"
	"github.com/OzkrOssa/freeradius-api/internal/core/domain"
	"github.com/OzkrOssa/freeradius-api/internal/core/port"
	"github.com/OzkrOssa/freeradius-api/internal/core/utils"
)

type AuthService struct {
	repo port.UserRepository
	ts   port.TokenService
}

func NewAuthService(repo port.UserRepository, ts port.TokenService) *AuthService {
	return &AuthService{repo, ts}
}

func (as *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := as.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.DataNotFoundError) {
			return "", domain.InvalidCredentialsError
		}
		return "", domain.InternalError
	}

	err = utils.ComparePassword(password, user.Password)
	if err != nil {
		return "", domain.InvalidCredentialsError
	}

	accessToken, err := as.ts.CreateToken(user)
	if err != nil {
		return "", domain.TokenCreationError
	}

	return accessToken, nil
}

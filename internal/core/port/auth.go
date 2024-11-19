package port

import (
	"context"
	"github.com/OzkrOssa/freeradius-api/internal/core/domain"
)

type TokenService interface {
	CreateToken(user *domain.User) (string, error)
	VerifyToken(token string) (*domain.TokenPayload, error)
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
}

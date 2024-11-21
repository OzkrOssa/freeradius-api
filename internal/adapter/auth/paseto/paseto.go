package paseto

import (
	"github.com/OzkrOssa/freeradius-api/internal/adapter/config"
	"github.com/OzkrOssa/freeradius-api/internal/core/domain"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"os"
	"strconv"
	"time"
)

type Token struct {
	token    *paseto.JSONToken
	key      *paseto.V2
	duration time.Duration
}

func New(config *config.Token) (*Token, error) {
	duration, err := time.ParseDuration(config.Duration)
	if err != nil {
		return nil, domain.TokenDurationError
	}

	token := &paseto.JSONToken{}
	key := paseto.NewV2()

	return &Token{
		token:    token,
		key:      key,
		duration: duration,
	}, nil
}

func (pt *Token) CreateToken(user *domain.User) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", domain.InternalError
	}

	payload := &domain.TokenPayload{
		ID:     id,
		UserID: user.ID,
	}

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(pt.duration)

	pt.token.IssuedAt = issuedAt
	pt.token.NotBefore = issuedAt
	pt.token.Expiration = expiredAt

	pt.token.Set("id", payload.ID.String())
	pt.token.Set("user_id", strconv.FormatUint(payload.UserID, 10))

	secretKey := []byte(os.Getenv("SECRET"))
	if len(secretKey) == 0 {
		return "", domain.InternalError
	}

	token, err := pt.key.Encrypt(secretKey, *pt.token, nil)
	if err != nil {
		return "", domain.InternalError
	}

	return token, nil
}

func (pt *Token) VerifyToken(token string) (*domain.TokenPayload, error) {
	var payload domain.TokenPayload

	secretKey := []byte(os.Getenv("SECRET"))
	if len(secretKey) == 0 {
		return nil, domain.InternalError
	}

	err := pt.key.Decrypt(token, secretKey, &payload, nil)
	if err != nil {
		return nil, domain.InvalidTokenError
	}

	if time.Now().After(pt.token.Expiration) {
		return nil, domain.ExpiredTokenError
	}

	return &payload, nil
}

package domain

import "github.com/google/uuid"

type TokenPayload struct {
	ID     uuid.UUID
	UserID uint64
}

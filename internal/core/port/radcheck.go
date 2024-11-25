package port

import (
	"context"
	"github.com/OzkrOssa/freeradius-api/internal/core/domain"
)

type RadCheckRepository interface {
	CreateRadCheck(ctx context.Context, radcheck *domain.RadCheck) (*domain.RadCheck, error)
	GetRadCheckByID(ctx context.Context, id uint64) (*domain.RadCheck, error)
	GetRadCheckByUserName(ctx context.Context, id uint64) (*domain.RadCheck, error)
	ListRadChecks(ctx context.Context, skip, limit uint64) ([]domain.RadCheck, error)
	UpdateRadCheck(ctx context.Context, radcheck *domain.RadCheck) (*domain.RadCheck, error)
	DeleteRadCheck(ctx context.Context, id uint64) error
}

type RadCheckService interface {
	CreateRadCheck(ctx context.Context, radcheck *domain.RadCheck) (*domain.RadCheck, error)
	GetRadCheck(ctx context.Context, id uint64) (*domain.RadCheck, error)
	ListRadChecks(ctx context.Context, skip, limit uint64) ([]domain.RadCheck, error)
	UpdateRadCheck(ctx context.Context, radcheck *domain.RadCheck) (*domain.RadCheck, error)
	DeleteRadCheck(ctx context.Context, id uint64) error
}

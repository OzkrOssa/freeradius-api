package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/OzkrOssa/freeradius-api/internal/adapter/storage/postgres"
	"github.com/OzkrOssa/freeradius-api/internal/core/domain"
)

type RadCheck struct {
	db *postgres.DB
}

func NewRadCheckRepository(db *postgres.DB) *RadCheck {

	return &RadCheck{db}
}

func (rd *RadCheck) CreateRadCheck(ctx context.Context, radcheck *domain.RadCheck) (*domain.RadCheck, error) {
	query := rd.db.Insert("radcheck").Columns("username", "attribute", "op", "value").Values(radcheck.UserName, radcheck.Attribute, radcheck.Operation, radcheck.Value).Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	err = rd.db.QueryRow(ctx, sql, args...).Scan(&radcheck.ID, &radcheck.UserName, &radcheck.Attribute, &radcheck.Operation, &radcheck.Value)
	if err != nil {
		if errorCode := rd.db.ErrorCode(err); errorCode == "23505" {
			return nil, domain.ConflictDataError
		}
		return nil, err
	}

	return radcheck, err
}

func (rd *RadCheck) GetRadCheckByID(ctx context.Context, id uint64) (*domain.RadCheck, error) {
	query := rd.db.Select("*").From("radcheck").Where("id = ?", id).Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var radCheck domain.RadCheck

	err = rd.db.QueryRow(ctx, sql, args...).Scan(&radCheck.ID, &radCheck.UserName, &radCheck.Attribute, &radCheck.Operation, &radCheck.Value)

	if err != nil {
		if errorCode := rd.db.ErrorCode(err); errorCode == "23505" {
			return nil, domain.ConflictDataError
		}
		return nil, err
	}

	return &radCheck, nil
}

func (rd *RadCheck) GetRadCheckByUserName(ctx context.Context, userName string) (*domain.RadCheck, error) {
	query := rd.db.Select("*").From("radcheck").Where("username = ?", userName).Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var radCheck domain.RadCheck

	err = rd.db.QueryRow(ctx, sql, args...).Scan(&radCheck.ID, &radCheck.UserName, &radCheck.Attribute, &radCheck.Operation, &radCheck.Value)

	if err != nil {
		if errorCode := rd.db.ErrorCode(err); errorCode == "23505" {
			return nil, domain.ConflictDataError
		}
		return nil, err
	}

	return &radCheck, nil
}

func (rd *RadCheck) ListRadChecks(ctx context.Context, skip, limit uint64) ([]domain.RadCheck, error) {
	var radcheck domain.RadCheck
	var listRadChecks []domain.RadCheck

	query := rd.db.Select("*").From("radcheck").OrderBy("id").Offset((skip - 1) * limit)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := rd.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&radcheck.ID,
			&radcheck.UserName,
			&radcheck.Attribute,
			&radcheck.Operation,
			&radcheck.Value,
		)
		if err != nil {
			return nil, err
		}
		listRadChecks = append(listRadChecks, radcheck)
	}

	return listRadChecks, err
}

func (rd *RadCheck) UpdateRadCheck(ctx context.Context, radcheck *domain.RadCheck) (*domain.RadCheck, error) {
	query := rd.db.Update("radcheck").Set("username", sq.Expr("COALESCE (?, username)", radcheck.UserName)).Set("attribute", sq.Expr("COALESCE (?, attribute)", radcheck.Attribute)).Set("op", sq.Expr("COALESCE (?, op)", radcheck.Operation)).Set("value", sq.Expr("COALESCE (?, value)", radcheck.Value)).Where(sq.Eq{"id": radcheck.ID}).Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	_, err = rd.db.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	err = rd.db.QueryRow(ctx, sql, args...).Scan(&radcheck.ID, &radcheck.UserName, &radcheck.Attribute, &radcheck.Operation, &radcheck.Value)
	if err != nil {
		if errCode := rd.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ConflictDataError
		}
		return nil, err
	}
	return radcheck, nil
}

func (rd *RadCheck) DeleteRadCheck(ctx context.Context, id uint64) error {
	query := rd.db.Delete("radcheck").Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = rd.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}

package repository

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/OzkrOssa/freeradius-api/internal/adapter/storage/postgres"
	"github.com/OzkrOssa/freeradius-api/internal/core/domain"
	"github.com/jackc/pgx/v5"
	"time"
)

type UserRepository struct {
	db *postgres.DB
}

func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {

	query := ur.db.Insert("users").Columns("name", "email", "password").Values(user.Name, user.Email, user.Password).Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errCode := ur.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ConflictDataError
		}
		return nil, err
	}

	return user, nil
}

func (ur UserRepository) GetUserById(ctx context.Context, id uint64) (*domain.User, error) {
	query := ur.db.Select("*").From("users").Where("id = ?", id).Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var user domain.User

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.DataNotFoundError
		}
		return nil, err
	}

	return &user, nil

}

func (ur UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := ur.db.Select("*").From("users").Where("email = ?", email).Limit(1)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	var user domain.User
	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.DataNotFoundError
		}
		return nil, err
	}
	return &user, nil
}

func (ur UserRepository) ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error) {
	var user domain.User
	var users []domain.User

	query := ur.db.Select("*").
		From("users").
		OrderBy("id").
		Limit(limit).
		Offset((skip - 1) * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := ur.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (ur UserRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {

	query := ur.db.Update("users").
		Set("name", sq.Expr("COALESCE(?, name)", user.Name)).
		Set("email", sq.Expr("COALESCE(?, email)", user.Email)).
		Set("password", sq.Expr("COALESCE(?, password)", user.Password)).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": user.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	_, err = ur.db.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	err = ur.db.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errCode := ur.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ConflictDataError
		}
		return nil, err
	}
	return user, nil
}

func (ur UserRepository) DeleteUser(ctx context.Context, id uint64) error {
	query := ur.db.Delete("users").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = ur.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

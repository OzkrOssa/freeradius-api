package postgres

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/OzkrOssa/freeradius-api/internal/adapter/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
	squirrel.StatementBuilderType
	url string
}

//go:embed migrations/*.sql
var migrationsFS embed.FS

func New(ctx context.Context, config *config.DB) (*DB, error) {
	url := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		config.Connection,
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return &DB{}, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return &DB{}, err
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	db := &DB{
		pool,
		psql,
		url,
	}

	return db, nil
}

func (db *DB) Migrate() error {
	driver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}

	migrations, err := migrate.NewWithSourceInstance("iofs", driver, db.url)
	if err != nil {
		return err
	}

	err = migrations.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
func (db *DB) ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	errors.As(err, &pgErr)
	return pgErr.Code
}

func (db *DB) Close() {
	db.Pool.Close()
}

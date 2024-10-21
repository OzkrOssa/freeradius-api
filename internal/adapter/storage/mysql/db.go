package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/OzkrOssa/freeradius-api/internal/adapter/config"
	"github.com/OzkrOssa/freeradius-api/internal/adapter/storage/mysql/sqlc"
	"github.com/go-sql-driver/mysql"
)

type DB struct {
	queries *sqlc.Queries
}

func New(ctx context.Context, config *config.DB) (*DB, error) {
	cfg := mysql.Config{
		User:   config.User,
		Passwd: config.Password,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%s", config.Url, config.Port),
		DBName: config.Name,
	}

	db, err := sql.Open(config.Connection, cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	pingError := db.Ping()
	if pingError != nil {
		db.Close()
		return nil, err
	}

	q := sqlc.New(db)

	return &DB{queries: q}, nil
}

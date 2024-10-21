package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/OzkrOssa/freeradius-api/internal/adapter/config"
	"github.com/OzkrOssa/freeradius-api/internal/adapter/storage/mysql"
	"github.com/OzkrOssa/freeradius-api/internal/adapter/storage/redis"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config, err := config.New()
	if err != nil {
		slog.Error("error reading config", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting the application", "app", config.App.Name, "env", config.App.Env)

	ctx := context.Background()
	db, err := mysql.New(ctx, config.DB)
	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}

	slog.Info("Successfully connected to the database", "db", config.DB.Connection)

	cache, err := redis.New(ctx, config.Redis)

	if err != nil {
		slog.Error("Error initializing cache connection", "error", err)
		os.Exit(1)
	}
	defer cache.Close()

	slog.Info("Successfully connected to the cache server")

}

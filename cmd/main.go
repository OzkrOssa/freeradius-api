package main

import (
	"context"
	"fmt"
	ps "github.com/OzkrOssa/freeradius-api/internal/adapter/auth/paseto"
	"github.com/OzkrOssa/freeradius-api/internal/adapter/config"
	"github.com/OzkrOssa/freeradius-api/internal/adapter/handler"
	"github.com/OzkrOssa/freeradius-api/internal/adapter/logger"
	"github.com/OzkrOssa/freeradius-api/internal/adapter/storage/postgres"
	"github.com/OzkrOssa/freeradius-api/internal/adapter/storage/postgres/repository"
	"github.com/OzkrOssa/freeradius-api/internal/adapter/storage/redis"
	"github.com/OzkrOssa/freeradius-api/internal/core/service"
	"log/slog"
	"os"
)

func main() {
	// Load environment variables
	c, err := config.New()
	if err != nil {
		slog.Error("error reading config", "err", err.Error())
		os.Exit(1)
	}

	// Set logger
	logger.Set(c.App)

	slog.Info("Starting the application", "app", c.App.Name, "env", c.App.Env)

	// Init database
	ctx := context.Background()
	db, err := postgres.New(ctx, c.DB)
	if err != nil {
		slog.Error("error initialize database", "err", err.Error())
		os.Exit(1)
	}

	defer db.Close()

	slog.Info("Successfully connected to the database", "db", c.DB.Connection)

	// Migrate database
	err = db.Migrate()
	if err != nil {
		slog.Error("Error migrating database", "error", err)
		os.Exit(1)
	}

	slog.Info("Successfully migrated the database")

	cache, err := redis.New(ctx, c.Redis)
	if err != nil {
		slog.Error("error initialize redis", "err", err.Error())
		os.Exit(1)
	}

	defer cache.Close()

	slog.Info("Successfully connected to the cache server")

	token, err := ps.New(c.Token)
	if err != nil {
		slog.Error("error initialize token service", "err", err.Error())
		os.Exit(1)
	}

	// Dependency injection
	// User
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cache)
	userHandler := handler.NewUserHandler(userService)

	//auth
	authService := service.NewAuthService(userRepo, token)
	authHandler := handler.NewAuthHandler(authService)

	// Init router
	r, err := handler.NewRouter(c.Http, token, authHandler, userHandler)
	if err != nil {
		slog.Error("error initialize router", "err", err.Error())
		os.Exit(1)
	}

	// Start server
	listenAddr := fmt.Sprintf("%s:%s", c.Http.Host, c.Http.Port)
	slog.Info("Starting the HTTP server", "listen_address", listenAddr)
	err = r.Run(listenAddr)
	if err != nil {
		slog.Error("error initialize http server", "err", err)
		os.Exit(1)
	}

}

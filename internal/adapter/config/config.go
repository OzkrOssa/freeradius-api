package config

import (
	"github.com/joho/godotenv"
	"os"
	"strings"
)

type (
	Container struct {
		App   *App
		Token *Token
		Redis *Redis
		DB    *DB
		Http  *Http
	}
	App struct {
		Name string
		Env  string
	}
	Token struct {
		Duration string
	}
	Redis struct {
		Address  string
		Port     string
		Password string
	}
	DB struct {
		Connection string
		Host       string
		Port       string
		Name       string
		User       string
		Password   string
	}
	Http struct {
		Host           string
		Port           string
		Env            string
		AllowedOrigins []string
	}
)

func New() (*Container, error) {
	if os.Getenv("APP_ENV") != "prod" {
		if err := godotenv.Load(); err != nil {
			return nil, err
		}
	}

	app := &App{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
	}
	token := &Token{
		Duration: os.Getenv("APP_TOKEN_DURATION"),
	}
	redis := &Redis{
		Address:  os.Getenv("REDIS_ADDRESS"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}
	db := &DB{
		Connection: os.Getenv("DB_CONNECTION"),
		Host:       os.Getenv("DB_HOST"),
		Port:       os.Getenv("DB_PORT"),
		Name:       os.Getenv("DB_NAME"),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
	}
	http := &Http{
		Host:           os.Getenv("HTTP_HOST"),
		Port:           os.Getenv("HTTP_PORT"),
		Env:            os.Getenv("HTTP_ENV"),
		AllowedOrigins: strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
	}
	return &Container{
		App:   app,
		Token: token,
		Redis: redis,
		DB:    db,
		Http:  http,
	}, nil
}

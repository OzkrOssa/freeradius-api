package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type (
	App struct {
		Name string
		Env  string
	}
	DB struct {
		Connection string
		Url        string
		Port       string
		Name       string
		User       string
		Password   string
	}
	Token struct {
		Duration string
	}
	Redis struct {
		Addr     string
		Password string
	}
	HTTP struct {
		Url            string
		Port           string
		Env            string
		AllowedOrigins string
	}
	Container struct {
		App   *App
		DB    *DB
		Token *Token
		Redis *Redis
		HTTP  *HTTP
	}
)

func New() (*Container, error) {
	if !strings.Contains("prod", os.Getenv("APP_ENV")) {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	app := &App{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
	}

	db := &DB{
		Connection: os.Getenv("DB_CONNECTION"),
		Url:        os.Getenv("DB_URL"),
		Port:       os.Getenv("DB_PORT"),
		Name:       os.Getenv("DB_NAME"),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
	}

	token := &Token{
		Duration: os.Getenv("TOKEN_DURATION"),
	}

	redis := &Redis{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	http := &HTTP{
		Url:            os.Getenv("HTTP_URL"),
		Port:           os.Getenv("HTTP_PORT"),
		Env:            os.Getenv("APP_ENV"),
		AllowedOrigins: os.Getenv("HTTP_ORIGIS"),
	}

	return &Container{
		App:   app,
		DB:    db,
		Token: token,
		Redis: redis,
		HTTP:  http,
	}, nil
}

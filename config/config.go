package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PROJECT_PORT string
	NGINX_URL    string

	DB_HOST        string
	DB_PORT        string
	DB_USERNAME    string
	DB_PASSWORD    string
	DB_AUTH_SOURCE string

	REDIS_HOST     string
	REDIS_PORT     string
	REDIS_PASSWORD string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("load env is failed")
	}

	return &Config{
		PROJECT_PORT: os.Getenv("PROJECT_PORT"),
		NGINX_URL:    os.Getenv("NGINX_URL"),

		DB_HOST:        os.Getenv("DB_HOST"),
		DB_PORT:        os.Getenv("DB_PORT"),
		DB_USERNAME:    os.Getenv("DB_USERNAME"),
		DB_PASSWORD:    os.Getenv("DB_PASSWORD"),
		DB_AUTH_SOURCE: os.Getenv("DB_AUTH_SOURCE"),

		REDIS_HOST:     os.Getenv("REDIS_HOST"),
		REDIS_PORT:     os.Getenv("REDIS_PORT"),
		REDIS_PASSWORD: os.Getenv("REDIS_PASSWORD"),
	}, nil
}

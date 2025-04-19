package data

import (
	"log"
	"ms-admin/config"
	"strconv"

	"github.com/gofiber/storage/redis/v3"
)

func ConnectRedis() (*redis.Storage, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	port, err := strconv.Atoi(cfg.REDIS_PORT)
	if err != nil {
		log.Fatal(err)
	}

	store := redis.New(redis.Config{
		Host:     cfg.REDIS_HOST,
		Port:     port,
		Password: cfg.REDIS_PASSWORD,
		Database: 0,
	})

	return store, nil
}

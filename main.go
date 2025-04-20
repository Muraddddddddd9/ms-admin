package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ms-admin/api/handlers"
	"ms-admin/config"
	"ms-admin/data"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/storage/redis/v3"
)

func AdminOnly(storage *redis.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session_id := c.Cookies("session_id")

		body, err := storage.Get(fmt.Sprintf("users:%v", session_id))
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Сессия не найдена"})
		}

		var user struct {
			Status string `json:"status"`
		}
		err = json.Unmarshal(body, &user)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Ошибка в получение данных"})
		}

		if user.Status != "admin" {
			return c.Status(404).JSON(fiber.Map{})
		}

		return c.Next()
	}
}

func main() {
	client, err := data.ConnectMongo()
	if err == nil {
		log.Println("Connect to MONGODB is succ")
	}
	defer client.Disconnect(context.TODO())

	storage, err := data.ConnectRedis()
	if err == nil {
		log.Println("Connect to REDIS is succ")
	}
	defer storage.Close()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Print(err)
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE",
		AllowCredentials: true,
	}))

	app.Post("/api/admin/create_data", AdminOnly(storage), func(c *fiber.Ctx) error {
		return handlers.CreateData(c, client)
	})

	app.Post("/api/admin/get_data/:collection", AdminOnly(storage), func(c *fiber.Ctx) error {
		return handlers.GetData(c, client)
	})

	app.Listen(cfg.PROJECT_PORT)
}

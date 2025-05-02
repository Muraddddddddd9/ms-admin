package main

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
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

		if user.Status != "админ" {
			return c.Status(404).JSON(fiber.Map{})
		}

		return c.Next()
	}
}

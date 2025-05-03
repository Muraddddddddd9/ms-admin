package main

import (
	"encoding/json"
	"fmt"
	"ms-admin/api/messages"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis/v3"
)

func AdminOnly(storage *redis.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session_id := c.Cookies("session_id")

		body, err := storage.Get(fmt.Sprintf("users:%v", session_id))
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"message": messages.ErrSessionNotFound})
		}

		var user struct {
			Status string `json:"status"`
		}
		err = json.Unmarshal(body, &user)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"message": messages.ErrGetData})
		}

		if user.Status != "админ" {
			return c.Status(404).JSON(fiber.Map{})
		}

		return c.Next()
	}
}

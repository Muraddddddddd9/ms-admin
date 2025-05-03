package main

import (
	"context"
	"encoding/json"
	"fmt"
	"ms-admin/api/messages"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func AdminOnly(storage *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session_id := c.Cookies("session_id")

		resGet, err := storage.Get(context.Background(), fmt.Sprintf("users:%v", session_id)).Result()
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"message": messages.ErrSessionNotFound})
		}

		var user struct {
			Status string `json:"status"`
		}
		err = json.Unmarshal([]byte(resGet), &user)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"message": messages.ErrGetData})
		}

		if user.Status != "админ" {
			return c.Status(404).JSON(fiber.Map{})
		}

		return c.Next()
	}
}

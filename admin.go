package main

import (
	"context"
	"encoding/json"
	"fmt"
	"ms-admin/api/constants"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

const (
	SessionKeyStart     = "session:%s"
	RedirectPathProfile = "/profile"
)

func AdminOnly(rdb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session")

		sessionKey := fmt.Sprintf(SessionKeyStart, session)

		userKey, err := rdb.Get(context.Background(), sessionKey).Result()
		if err == redis.Nil || userKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message":  constants.ErrSessionNotFound,
				"redirect": RedirectPathProfile,
			})
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message":  constants.ErrServerError,
				"redirect": RedirectPathProfile,
			})
		}

		userData, err := rdb.Get(context.Background(), userKey).Bytes()
		if err == redis.Nil || userData == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message":  constants.ErrUserNotFound,
				"redirect": RedirectPathProfile,
			})
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message":  constants.ErrServerError,
				"redirect": RedirectPathProfile,
			})
		}

		var user struct {
			Status string `json:"status"`
		}
		err = json.Unmarshal(userData, &user)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"message":  constants.ErrGetData,
				"redirect": RedirectPathProfile,
			})
		}

		if user.Status != "админ" {
			return c.Status(301).JSON(fiber.Map{
				"redirect": RedirectPathProfile,
			})
		}

		return c.Next()
	}
}

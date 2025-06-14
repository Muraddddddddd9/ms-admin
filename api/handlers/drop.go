package handlers

import (
	"context"
	"fmt"
	"ms-admin/api/constants"
	"ms-admin/api/services"
	loconfig "ms-admin/config"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Drop(c *fiber.Ctx, db *mongo.Database, cfg *loconfig.LocalConfig) error {
	var collectionsData CollectionStruct

	if err := c.BodyParser(&collectionsData); err != nil {
		services.Logging(db, "/api/admin/drop", c.Method(), strconv.Itoa(fiber.StatusBadRequest), collectionsData, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidInput,
		})
	}

	if len(collectionsData.Collections) == 0 {
		services.Logging(db, "/api/admin/drop", c.Method(), strconv.Itoa(fiber.StatusConflict), collectionsData, constants.ErrInputCollection)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": constants.ErrInputCollection,
		})
	}

	var errMsg []string

	for _, collection := range collectionsData.Collections {
		if collection == constants.TeacherCollection {
			_, err := db.Collection(collection).
				DeleteMany(
					context.Background(),
					bson.M{
						"status": bson.M{"$ne": constants.AdminStatus},
						"email":  bson.M{"$ne": cfg.ADMIN_EMAIL},
					},
				)
			if err != nil {
				errMsg = append(errMsg, fmt.Sprintf("%v: %v", constants.ErrDeleteData, constants.TeacherCollection))
			}
			continue
		}
		if collection == constants.StatusCollection {
			continue
		}
		err := db.Collection(collection).Drop(context.Background())
		if err != nil {
			errMsg = append(errMsg, fmt.Sprintf("%v: %v", constants.ErrDeleteData, collection))
		}
	}

	if len(errMsg) != 0 {
		services.Logging(db, "/api/admin/drop", c.Method(), strconv.Itoa(fiber.StatusBadRequest), collectionsData, errMsg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errMsg,
		})
	}

	services.Logging(db, "/api/admin/drop", c.Method(), strconv.Itoa(fiber.StatusAccepted), collectionsData, constants.SuccDeleteCollection)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccDeleteCollection,
	})
}

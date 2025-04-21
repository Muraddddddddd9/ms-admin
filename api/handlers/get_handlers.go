package handlers

import (
	"ms-admin/api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetData(c *fiber.Ctx, db *mongo.Database) error {
	collection := c.Params("collection")

	var data interface{}
	var err error

	switch collection {
	case "students":
		data, err = services.ReadStudent(db)
	case "teachers":
		data, err = services.ReadTeachers(db)
	case "groups":
		data, err = services.ReadGroups(db)
	case "objects":
		data, err = services.ReadObjects(db)
	case "objects_groups":
		data, err = services.ReadObjectsGroups(db)
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"data": data,
	})
}

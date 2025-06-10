package handlers

import (
	"fmt"
	"ms-admin/api/constants"
	"ms-admin/api/services"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type getFunc func(db *mongo.Database, page, pageSize int) (map[string]interface{}, error)

var handlersMapGet = map[string]getFunc{
	constants.StudentCollection:     services.ReadStudents,
	constants.TeacherCollection:     services.ReadTeachers,
	constants.GroupCollection:       services.ReadGroups,
	constants.ObjectCollection:      services.ReadObjects,
	constants.ObjectGroupCollection: services.ReadObjectsGroups,
	constants.StatusCollection:      services.ReadStatuses,
}

func GetData(c *fiber.Ctx, db *mongo.Database) error {
	collection := c.Params("collection")
	page := c.QueryInt("page")
	pageSize := c.QueryInt("pageSize")

	handler, exists := handlersMapGet[strings.TrimSpace(collection)]
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrCollectionNotFound,
		})
	}

	var err error
	data, err := handler(db, page, pageSize)

	if err != nil {
		services.Logging(db, fmt.Sprintf("/api/admin/get_data/%v", collection), "GET", "400", data, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"data": data,
	})
}

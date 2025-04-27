package handlers

import (
	"fmt"
	"ms-admin/api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetData(c *fiber.Ctx, db *mongo.Database) error {
	collection := c.Params("collection")

	var data interface{}
	var selectData interface{}
	var header []string
	var err error

	switch collection {
	case "students":
		data, header, selectData, err = services.ReadStudents(db)
	case "teachers":
		data, header, selectData, err = services.ReadTeachers(db)
	case "groups":
		data, header, selectData, err = services.ReadGroups(db)
	case "objects":
		data, header, err = services.ReadObjects(db)
	case "objects_groups":
		data, header, selectData, err = services.ReadObjectsGroups(db)
	case "statuses":
		data, header, err = services.ReadStatuses(db)
	case "logs":
		data, header, err = services.ReadLogs(db)
	}

	if err != nil {
		services.Logging(db, fmt.Sprintf("/api/admin/get_data/%v", collection), "GET", "400", data, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"data":       data,
		"header":     header,
		"selectData": selectData,
	})
}

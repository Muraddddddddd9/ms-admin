package handlers

import (
	"fmt"
	"ms-admin/api/constants"
	"ms-admin/api/services"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetData(c *fiber.Ctx, db *mongo.Database) error {
	collection := c.Params("collection")

	var data interface{}
	var selectData interface{}
	var header []string
	var err error

	switch strings.TrimSpace(collection) {
	case constants.StudentCollection:
		data, header, selectData, err = services.ReadStudents(db)
	case constants.TeacherCollection:
		data, header, selectData, err = services.ReadTeachers(db)
	case constants.GroupCollection:
		data, header, selectData, err = services.ReadGroups(db)
	case constants.ObjectCollection:
		data, header, err = services.ReadObjects(db)
	case constants.ObjectGroupCollection:
		data, header, selectData, err = services.ReadObjectsGroups(db)
	case constants.StatusCollection:
		data, header, err = services.ReadStatuses(db)
	case constants.LogsCollection:
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

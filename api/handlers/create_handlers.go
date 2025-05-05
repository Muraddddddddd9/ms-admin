package handlers

import (
	"encoding/json"
	"fmt"
	"ms-admin/api/constants"
	"ms-admin/api/services"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateDataStruct struct {
	Collection string          `json:"collection"`
	NewData    json.RawMessage `json:"new_data"`
}

func CreateData(c *fiber.Ctx, db *mongo.Database) error {
	var data CreateDataStruct
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidInput,
		})
	}

	var newId interface{}
	var err error
	switch strings.TrimSpace(data.Collection) {
	case constants.StudentCollection:
		newId, err = services.CreateStudents(db, data.NewData)
	case constants.TeacherCollection:
		newId, err = services.CreateTeachers(db, data.NewData)
	case constants.GroupCollection:
		newId, err = services.CreateGroups(db, data.NewData)
	case constants.ObjectCollection:
		newId, err = services.CreateObjects(db, data.NewData)
	case constants.ObjectGroupCollection:
		newId, err = services.CreateObjectsGroups(db, data.NewData)
	case constants.StatusCollection:
		newId, err = services.CreateStatuses(db, data.NewData)
	}

	var dataForLog any
	_ = json.Unmarshal(data.NewData, &dataForLog)

	if err != nil {
		services.Logging(db, "/api/admin/create_data", "POST", "400", dataForLog, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	services.Logging(db, "/api/admin/create_data", "POST", "202", dataForLog, nil)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": fmt.Sprintf(constants.SuccDataAdd, newId.(primitive.ObjectID).Hex()),
	})
}

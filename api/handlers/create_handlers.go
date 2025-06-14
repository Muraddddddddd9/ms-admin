package handlers

import (
	"encoding/json"
	"fmt"
	"ms-admin/api/constants"
	"ms-admin/api/services"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateDataStruct struct {
	Collection string          `json:"collection"`
	NewData    json.RawMessage `json:"new_data"`
}

type createFunc func(db *mongo.Database, data json.RawMessage) (interface{}, error)

var handlersMapCreate = map[string]createFunc{
	constants.StudentCollection:     services.CreateStudents,
	constants.TeacherCollection:     services.CreateTeachers,
	constants.GroupCollection:       services.CreateGroups,
	constants.ObjectCollection:      services.CreateObjects,
	constants.ObjectGroupCollection: services.CreateObjectsGroups,
	constants.StatusCollection:      services.CreateStatuses,
}

func CreateData(c *fiber.Ctx, db *mongo.Database) error {
	var data CreateDataStruct
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidInput,
		})
	}

	handler, exists := handlersMapCreate[strings.TrimSpace(data.Collection)]
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrCollectionNotFound,
		})
	}

	newId, err := handler(db, data.NewData)

	var dataForLog any
	_ = json.Unmarshal(data.NewData, &dataForLog)

	if err != nil {
		services.Logging(db, "/api/admin/create_data", c.Method(), strconv.Itoa(fiber.StatusBadRequest), dataForLog, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	services.Logging(db, "/api/admin/create_data", c.Method(), strconv.Itoa(fiber.StatusAccepted), dataForLog, nil)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": fmt.Sprintf(constants.SuccDataAdd, newId.(primitive.ObjectID).Hex()),
	})
}

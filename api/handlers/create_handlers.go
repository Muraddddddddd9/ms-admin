package handlers

import (
	"encoding/json"
	"fmt"
	"ms-admin/api/services"

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
			"message": "Данные введены не верно",
		})
	}

	var newId interface{}
	var err error
	switch data.Collection {
	case "students":
		newId, err = services.CreateStudents(db, data.NewData)
	case "teachers":
		newId, err = services.CreateTeachers(db, data.NewData)
	case "groups":
		newId, err = services.CreateGroups(db, data.NewData)
	case "objects":
		newId, err = services.CreateObjects(db, data.NewData)
	case "objects_groups":
		newId, err = services.CreateObjectsGroups(db, data.NewData)
	case "statuses":
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
		"message": fmt.Sprintf("Данные добавлены с ID: %v", newId.(primitive.ObjectID).Hex()),
	})
}

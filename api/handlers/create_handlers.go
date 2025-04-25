package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

	var insertData any
	var err error

	switch data.Collection {
	case "students":
		insertData, err = services.CreateStudents(db, data.NewData)
	case "teachers":
		insertData, err = services.CreateTeachers(db, data.NewData)
	case "groups":
		insertData, err = services.CreateGroups(db, data.NewData)
	case "objects":
		insertData, err = services.CreateObjects(db, data.NewData)
	case "objects_groups":
		insertData, err = services.CreateObjectsGroups(db, data.NewData)
	case "statuses":
		insertData, err = services.CreateStatuses(db, data.NewData)
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	collectionID, err := db.Collection(data.Collection).InsertOne(context.TODO(), insertData)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Данные не были добавлены",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": fmt.Sprintf("Данные добавлены с ID: %v", collectionID.InsertedID.(primitive.ObjectID).Hex()),
	})
}

package handlers

import (
	"ms-admin/api/constants"
	"ms-admin/api/services"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UpdateDataStruct struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	Collection string             `bson:"collection" json:"collection"`
	Label      string             `bson:"label" json:"label"`
	NewData    string             `bson:"new_data" json:"new_data"`
}

type updateFunc func(
	db *mongo.Database,
	rdb *redis.Client,
	collection string,
	id primitive.ObjectID,
	label string,
	newData string,
) error

var handlersMapUpdate = map[string]updateFunc{
	constants.StudentCollection:     services.UpdateStudents,
	constants.TeacherCollection:     services.UpdateTeachers,
	constants.GroupCollection:       services.UpdateGroups,
	constants.ObjectCollection:      services.UpdateObjects,
	constants.ObjectGroupCollection: services.UpdateObjectsGroups,
	constants.StatusCollection:      services.UpdateStatuses,
}

func UpdateData(c *fiber.Ctx, db *mongo.Database, rdb *redis.Client) error {
	var data UpdateDataStruct
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidInput,
		})
	}

	data.Collection = strings.TrimSpace(strings.ToLower(data.Collection))
	data.Label = strings.TrimSpace(strings.ToLower(data.Label))
	data.NewData = strings.TrimSpace(data.NewData)

	handler, exists := handlersMapUpdate[data.Collection]
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrCollectionNotFound,
		})
	}

	err := handler(db, rdb, data.Collection, data.ID, data.Label, data.NewData)
	if err != nil {
		services.Logging(db, "/api/admin/update_data", "POST", "404", data, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	services.Logging(db, "/api/admin/update_data", "POST", "202", data, nil)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccDataUpdate,
	})
}

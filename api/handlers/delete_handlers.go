package handlers

import (
	"ms-admin/api/constants"
	"ms-admin/api/services"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeleteDataStruct struct {
	ID         []primitive.ObjectID `bson:"_id" json:"_id"`
	Collection string               `bson:"collection" json:"collection"`
}

type deleteFunc func(
	db *mongo.Database,
	collectionName string,
	ids []primitive.ObjectID,
) (string, error)

var handlersMapDelete = map[string]deleteFunc{
	constants.StudentCollection:     services.DeleteStudents,
	constants.TeacherCollection:     services.DeleteTeachers,
	constants.GroupCollection:       services.DeleteGroups,
	constants.ObjectCollection:      services.DeleteObjects,
	constants.ObjectGroupCollection: services.DeleteObjectsGroups,
	constants.StatusCollection:      services.DeleteStatuses,
}

func DeleteData(c *fiber.Ctx, db *mongo.Database) error {
	var deleteData DeleteDataStruct
	if err := c.BodyParser(&deleteData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidInput,
		})
	}

	deleteData.Collection = strings.TrimSpace(strings.ToLower(deleteData.Collection))

	handler, exists := handlersMapDelete[strings.TrimSpace(deleteData.Collection)]
	if !exists {
		services.Logging(db, "/api/admin/delete_data", c.Method(), strconv.Itoa(fiber.StatusBadRequest), deleteData, constants.ErrCollectionNotFound)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrCollectionNotFound,
		})
	}

	res, err := handler(db, deleteData.Collection, deleteData.ID)
	if err != nil {
		services.Logging(db, "/api/admin/delete_data", c.Method(), strconv.Itoa(fiber.StatusConflict), deleteData, err.Error())
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	services.Logging(db, "/api/admin/delete_data", c.Method(), strconv.Itoa(fiber.StatusAccepted), deleteData, nil)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": res,
	})
}

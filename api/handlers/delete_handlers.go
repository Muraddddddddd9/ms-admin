package handlers

import (
	"context"
	"fmt"
	"ms-admin/api/constants"
	"ms-admin/api/services"
	"ms-admin/api/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteData(c *fiber.Ctx, db *mongo.Database) error {
	var deleteData struct {
		ID         []primitive.ObjectID `bson:"_id" json:"_id"`
		Collection string               `bson:"collection" json:"collection"`
	}

	if err := c.BodyParser(&deleteData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidInput,
		})
	}

	deleteData.Collection = strings.TrimSpace(strings.ToLower(deleteData.Collection))

	var countNoneDelete int = 0
	var errFind error
	var errResult error

	var collectionDependencies = map[string][]struct {
		collection string
		field      string
	}{
		constants.GroupCollection: {
			{constants.StudentCollection, "group"},
			{constants.ObjectGroupCollection, "group"},
		},
		constants.TeacherCollection: {
			{constants.GroupCollection, "teacher"},
			{constants.ObjectGroupCollection, "teacher"},
		},
		constants.StatusCollection: {
			{constants.StudentCollection, "status"},
			{constants.TeacherCollection, "status"},
		},
		constants.ObjectCollection: {
			{constants.ObjectGroupCollection, "object"},
		},
	}

	for _, v := range deleteData.ID {
		filter := bson.M{"_id": v}

		if deps, ok := collectionDependencies[deleteData.Collection]; ok {
			for _, dep := range deps {
				if errFind = utils.CheckDataOtherTable(db, dep.collection, bson.M{dep.field: v}); errFind != nil {
					countNoneDelete++
					errResult = errFind
					continue
				}
			}
			if errFind != nil {
				continue
			}
		}

		_, err := db.Collection(deleteData.Collection).DeleteOne(context.Background(), filter)
		if err != nil {
			services.Logging(db, "/api/admin/delete_data", "POST", "400", deleteData, err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": constants.ErrDeleteData,
			})
		}
	}

	if countNoneDelete >= 1 {
		services.Logging(db, "/api/admin/delete_data", "POST", "409", deleteData, errResult.Error())
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": fmt.Sprint(errResult),
		})
	}

	services.Logging(db, "/api/admin/delete_data", "POST", "202", deleteData, nil)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": fmt.Sprintf(constants.SuccDataDelete, len(deleteData.ID)-countNoneDelete, len(deleteData.ID)),
	})
}

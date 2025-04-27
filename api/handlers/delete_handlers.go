package handlers

import (
	"context"
	"fmt"
	"ms-admin/api/services"

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
			"message": "Данные введены не верно",
		})
	}

	var countNoneDelete int = 0
	var errFind error
	var errResult error

	for _, v := range deleteData.ID {
		filter := bson.M{
			"_id": v,
		}

		switch deleteData.Collection {
		case "groups":
			errFind = services.CheckDataOtherTable(db, "students", bson.M{"group": v})
			if errFind != nil {
				countNoneDelete++
				errResult = errFind
				continue
			}
			errFind = services.CheckDataOtherTable(db, "objects_groups", bson.M{"group": v})
			if errFind != nil {
				countNoneDelete++
				errResult = errFind
				continue
			}
		case "teachers":
			errFind = services.CheckDataOtherTable(db, "groups", bson.M{"teacher": v})
			if errFind != nil {
				countNoneDelete++
				errResult = errFind
				continue
			}
			errFind = services.CheckDataOtherTable(db, "objects_groups", bson.M{"teacher": v})
			if errFind != nil {
				countNoneDelete++
				errResult = errFind
				continue
			}
		case "statuses":
			errFind = services.CheckDataOtherTable(db, "students", bson.M{"status": v})
			if errFind != nil {
				countNoneDelete++
				errResult = errFind
				continue
			}
			errFind = services.CheckDataOtherTable(db, "teachers", bson.M{"status": v})
			if errFind != nil {
				countNoneDelete++
				errResult = errFind
				continue
			}
		case "objects":
			errFind = services.CheckDataOtherTable(db, "objects_groups", bson.M{"object": v})
			if errFind != nil {
				countNoneDelete++
				errResult = errFind
				continue
			}
		}

		_, err := db.Collection(deleteData.Collection).DeleteOne(context.TODO(), filter)
		if err != nil {
			services.Logging(db, "/api/admin/delete_data", "POST", "400", deleteData, err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Ошибка в удаление",
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
		"message": fmt.Sprintf("Было удалено %v из %v", len(deleteData.ID)-countNoneDelete, len(deleteData.ID)),
	})
}

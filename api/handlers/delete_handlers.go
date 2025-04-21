package handlers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteData(c *fiber.Ctx, db *mongo.Database) error {
	var deleteData struct {
		ID         []primitive.ObjectID `bson:"_id" json:"_id"`
		Collection string               `bson:"collection" json:"collection"`
		Test       []string             `bson:"test" json:"test"`
	}

	if err := c.BodyParser(&deleteData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Данные введены не верно",
		})
	}

	fmt.Print(deleteData)
	for _, v := range deleteData.ID {
		filter := bson.M{
			"_id": v,
		}
		result, err := db.Collection(deleteData.Collection).DeleteOne(context.TODO(), filter)
		fmt.Print(result, err)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Ошибка в удаление",
			})
		}
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"messsage": "Данные были удалены",
	})
}

package handlers

import (
	"context"
	"ms-admin/api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateData(c *fiber.Ctx, db *mongo.Database) error {
	var data struct {
		ID         primitive.ObjectID `bson:"_id" json:"_id"`
		Collection string             `bson:"collection" json:"collection"`
		Label      string             `bson:"label" json:"label"`
		NewData    string             `bson:"new_data" json:"new_data"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Данные введены не верно",
		})
	}

	var err error
	switch data.Collection {
	case "students", "teachers":
		if data.Label == "email" || data.Label == "telegram" {
			err = services.CheckReplica(db, data.Collection, bson.M{data.Label: data.NewData})
		}
	case "statuses":
		if data.Label == "status" {
			err = services.CheckReplica(db, data.Collection, bson.M{data.Label: data.NewData})
		}
	case "objects":
		if data.Label == "object" {
			err = services.CheckReplica(db, data.Collection, bson.M{data.Label: data.NewData})
		}
	case "groups":
		if data.Label == "group" {
			err = services.CheckReplica(db, data.Collection, bson.M{data.Label: data.NewData})
		}
	}

	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	filter := bson.M{"_id": data.ID}
	update := bson.M{
		"$set": bson.M{
			data.Label: data.NewData,
		},
	}

	_, err = db.Collection(data.Collection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Обновленние данных провалилась",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Данные обновлены",
	})
}

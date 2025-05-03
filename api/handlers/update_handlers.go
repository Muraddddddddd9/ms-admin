package handlers

import (
	"context"
	"ms-admin/api/messages"
	"ms-admin/api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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
			"message": messages.ErrInvalidInput,
		})
	}

	var err error

	switch data.Collection {
	case "students", "teachers":
		if data.Label == "email" || data.Label == "telegram" {
			err = services.CheckReplica(db, data.Collection, bson.M{data.Label: data.NewData})
		}

		if data.Label == "password" {
			bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(data.NewData), bcrypt.DefaultCost)
			data.NewData = string(bcryptPassword)
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
		services.Logging(db, "/api/admin/update_data", "POST", "409", data, err.Error())
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	filter := bson.M{"_id": data.ID}
	var update bson.M

	if data.Label == "group" && data.Collection != "groups" ||
		data.Label == "status" && data.Collection != "statuses" ||
		data.Label == "object" && data.Collection == "objects_groups" ||
		data.Label == "teacher" {
		newData, _ := primitive.ObjectIDFromHex(data.NewData)
		update = bson.M{
			"$set": bson.M{
				data.Label: newData,
			},
		}
	} else {
		update = bson.M{
			"$set": bson.M{
				data.Label: data.NewData,
			},
		}
	}

	_, err = db.Collection(data.Collection).UpdateOne(context.Background(), filter, update)
	if err != nil {
		services.Logging(db, "/api/admin/update_data", "POST", "404", data, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": messages.ErrUpdateData,
		})
	}

	services.Logging(db, "/api/admin/update_data", "POST", "202", data, nil)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": messages.SuccDataUpdate,
	})
}

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Data struct {
	Collection string          `json:"collection"`
	NewData    json.RawMessage `json:"new_data"`
}


func CreateTeachers(db *mongo.Database, data Data, insertData interface{}) error {
	var teacher Teachers
	if err := json.Unmarshal(data.NewData, &teacher); err != nil {
		return fmt.Errorf("%v", "Неверные данные учителя")
	}

	err := checkReplica(db, data, bson.M{"email": teacher.Email})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	insertData = teacher
	return nil
}

func CreateGroups(db *mongo.Database, data Data, insertData interface{}) error {
	var group Groups
	if err := json.Unmarshal(data.NewData, &group); err != nil {
		return fmt.Errorf("%v", "Неверные данные группы")
	}

	err := checkReplica(db, data, bson.M{"group": group.Group})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	group.TeachersId, err = primitive.ObjectIDFromHex(group.TeachersId.String())
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	var findTeacher any
	err = db.Collection("teachers").FindOne(context.TODO(), bson.M{"_id": group.TeachersId}).Decode(&findTeacher)
	if err != nil {
		return fmt.Errorf("%v", "Учитель не найден")
	}

	insertData = group
	return nil
}

func CreateData(c *fiber.Ctx, client *mongo.Client) error {
	var data Data
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Данные введены не верно",
		})
	}

	db := client.Database("diary")
	var insertData any

	switch data.Collection {
	case "students":
		err := CreateStudent(db, data, &insertData)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	case "teachers":
		err := CreateTeachers(db, data, &insertData)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	case "groups":
		err := CreateGroups(db, data, &insertData)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	case "objects":
		var object Objects
		if err := json.Unmarshal(data.NewData, &object); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Неверные данные предмета",
			})
		}
		insertData = object
		err := checkReplica(db, data, bson.M{"object": object.Object})
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	case "objects_for_groups":
		var objects_for_groups ObjectsForGroups
		if err := json.Unmarshal(data.NewData, &objects_for_groups); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Неверные данные предмета для группы",
			})
		}
		insertData = objects_for_groups
	}

	collectionID, err := db.Collection(data.Collection).InsertOne(context.TODO(), insertData)
	if err != nil {
		log.Print(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Данные не были добавлены",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": fmt.Sprintf("Данные добавлены с ID: %v", collectionID),
	})
	return nil
}

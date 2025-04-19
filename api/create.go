package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Data struct {
	Collection string          `json:"collection"`
	NewData    json.RawMessage `json:"new_data"`
}

type Students struct {
	Name       string   `bson:"name"`
	Surname    string   `bson:"surname"`
	Patronymic string   `bson:"patronymic,omitempty"`
	Group      string   `bson:"group"`
	Email      string   `bson:"email"`
	Password   string   `bson:"password"`
	Telegram   string   `bson:"telegram,omitempty"`
	Diplomas   []string `bson:"diplomas,omitempty"`
	Ips        []string `bson:"ips,omitempty"`
	Status     string   `bson:"status"`
}

type Teachers struct {
	Name       string   `bson:"name"`
	Surname    string   `bson:"surname"`
	Patronymic string   `bson:"patronymic,omitempty"`
	Email      string   `bson:"email"`
	Password   string   `bson:"password"`
	Telegram   string   `bson:"telegram,omitempty"`
	Ips        []string `bson:"ips,omitempty"`
	Status     string   `bson:"status"`
}

type Groups struct {
	Name       string `bson:"name"`
	TeachersId string `bson:"teachers_id"`
}

type Objects struct {
	Name string `bson:"name"`
}

type ObjectsForGroups struct {
	NameId     string `json:"name_id"`
	GroupId    string `json:"group_id"`
	TeachersId string `json:"teachers_id"`
}

func CreateData(c *fiber.Ctx, client *mongo.Client) error {
	var data Data
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Данные введены не верно",
		})
	}

	var insertData interface{}
	var email string
	var group_num string
	var object_name string

	switch data.Collection {
	case "students":
		var student Students
		if err := json.Unmarshal(data.NewData, &student); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Неверные данные студента",
			})
		}
		insertData = student
		email = student.Email
	case "teachers":
		var teacher Teachers
		if err := json.Unmarshal(data.NewData, &teacher); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Неверные данные учителя",
			})
		}
		insertData = teacher
		email = teacher.Email
	case "groups":
		var group Groups
		if err := json.Unmarshal(data.NewData, &group); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Неверные данные группы",
			})
		}
		insertData = group
		group_num = group.Name
	case "objects":
		var object Objects
		if err := json.Unmarshal(data.NewData, &object); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Неверные данные предмета",
			})
		}
		insertData = object
		object_name = object.Name
	case "objects_for_groups":
		var objects_for_groups ObjectsForGroups
		if err := json.Unmarshal(data.NewData, &objects_for_groups); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Неверные данные предмета для группы",
			})
		}
		insertData = objects_for_groups
	}

	db := client.Database("diary")
	if data.Collection == "students" || data.Collection == "teachers" {
		filter := bson.M{
			"email": email,
		}
		err := db.Collection(data.Collection).FindOne(context.TODO(), filter)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Email уже существует",
			})
		}
	}

	if data.Collection == "groups" {
		filter := bson.M{
			"name": group_num,
		}
		err := db.Collection(data.Collection).FindOne(context.TODO(), filter)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Группа уже существует",
			})
		}
	}

	if data.Collection == "objects" {
		filter := bson.M{
			"name": object_name,
		}
		err := db.Collection(data.Collection).FindOne(context.TODO(), filter)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Предмет уже существует",
			})
		}
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
}

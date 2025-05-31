package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"ms-admin/api/constants"
	"ms-admin/api/services"
	"ms-admin/api/utils"
	"strings"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func UpdateData(c *fiber.Ctx, db *mongo.Database, rdb *redis.Client) error {
	var data struct {
		ID         primitive.ObjectID `bson:"_id" json:"_id"`
		Collection string             `bson:"collection" json:"collection"`
		Label      string             `bson:"label" json:"label"`
		NewData    string             `bson:"new_data" json:"new_data"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidInput,
		})
	}

	data.Collection = strings.TrimSpace(strings.ToLower(data.Collection))
	data.Label = strings.TrimSpace(strings.ToLower(data.Label))
	data.NewData = strings.TrimSpace(data.NewData)

	var err error

	switch data.Collection {
	case constants.StudentCollection, constants.TeacherCollection:
		if data.Label == "email" || data.Label == "telegram" {
			err = utils.CheckReplica(db, data.Collection, bson.M{data.Label: data.NewData})
		}

		if data.Label == "password" {
			bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(data.NewData), bcrypt.DefaultCost)
			data.NewData = string(bcryptPassword)
		}
	case constants.StatusCollection:
		if data.Label == "status" {
			err = utils.CheckReplica(db, data.Collection, bson.M{data.Label: data.NewData})
		}
	case constants.ObjectCollection:
		if data.Label == "object" {
			err = utils.CheckReplica(db, data.Collection, bson.M{data.Label: data.NewData})
		}
	case constants.GroupCollection:
		if data.Label == "group" {
			err = utils.CheckReplica(db, data.Collection, bson.M{data.Label: data.NewData})
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

	if data.Label == "group" && data.Collection != constants.GroupCollection ||
		data.Label == "status" && data.Collection != constants.StatusCollection ||
		data.Label == "object" && data.Collection == constants.ObjectGroupCollection ||
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
			"message": constants.ErrUpdateData,
		})
	}

	if data.Collection == constants.StudentCollection || data.Collection == constants.TeacherCollection {
		if resSession, _ := rdb.Get(context.Background(), fmt.Sprintf(constants.UserKeyStart, data.ID)).Result(); resSession != "" {
			var userData interface{}
			var userID string

			switch data.Collection {
			case constants.StudentCollection:
				studentRepo := mongodb.NewRepository[models.StudentsModel, struct{}](db.Collection(constants.StudentCollection))
				student, _ := studentRepo.FindOne(context.Background(), filter)

				groupRepo := mongodb.NewRepository[models.GroupsModel, struct{}](db.Collection(constants.GroupCollection))
				group, _ := groupRepo.FindOne(context.Background(), bson.M{"_id": student.Group})

				statusRepo := mongodb.NewRepository[models.StatusesModel, struct{}](db.Collection(constants.StatusCollection))
				status, _ := statusRepo.FindOne(context.Background(), bson.M{"_id": student.Status})

				userData, _ = json.Marshal(models.StudentsWithGroupAndStatusModel{
					ID:         student.ID,
					Name:       student.Name,
					Surname:    student.Surname,
					Patronymic: student.Patronymic,
					Group:      group.Group,
					Email:      student.Email,
					Password:   "",
					Telegram:   student.Telegram,
					Diplomas:   student.Diplomas,
					IPs:        student.IPs,
					Status:     status.Status,
				})
				userID = student.ID.Hex()
			case constants.TeacherCollection:
				teacherRepo := mongodb.NewRepository[models.TeachersModel, struct{}](db.Collection(constants.TeacherCollection))
				teacher, _ := teacherRepo.FindOne(context.Background(), filter)

				statusRepo := mongodb.NewRepository[models.StatusesModel, struct{}](db.Collection(constants.StatusCollection))
				status, _ := statusRepo.FindOne(context.Background(), bson.M{"_id": teacher.Status})

				userData, _ = json.Marshal(models.TeachersWithStatusModel{
					ID:         teacher.ID,
					Name:       teacher.Name,
					Surname:    teacher.Surname,
					Patronymic: teacher.Patronymic,
					Email:      teacher.Email,
					Password:   "",
					Telegram:   teacher.Telegram,
					IPs:        teacher.IPs,
					Status:     status.Status,
				})
				userID = teacher.ID.Hex()
			}

			rdb.Set(context.Background(), fmt.Sprintf(constants.UserKeyStart, userID), userData, redis.KeepTTL)
		}
	}

	services.Logging(db, "/api/admin/update_data", "POST", "202", data, nil)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccDataUpdate,
	})
}

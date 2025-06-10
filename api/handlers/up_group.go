package handlers

import (
	"context"
	"fmt"
	"ms-admin/api/constants"
	"ms-admin/api/services"
	"strconv"
	"strings"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckLastGroup(group int) bool {
	var lastGroup = []int{31, 32, 33, 34, 35, 36, 47, 48, 49}
	for _, gr := range lastGroup {
		if gr == group {
			return true
		}
	}

	return false
}

func DeleteUserData(db *mongo.Database, group primitive.ObjectID) error {
	objectGroupRepo := mongodb.NewRepository[models.ObjectsGroupsModel, struct{}](db.Collection(constants.ObjectGroupCollection))
	objectFindAll, err := objectGroupRepo.FindAll(context.Background(), bson.M{"group": group})
	if err != nil {
		return fmt.Errorf(constants.ErrObjectNotFound)
	}

	evaluationRepo := mongodb.NewRepository[models.EvaluationModel, struct{}](db.Collection(constants.EvaluationCollection))
	for _, object := range objectFindAll {
		err := evaluationRepo.DeleteMany(context.Background(), bson.M{"object": object.ID})
		if err != nil {
			fmt.Errorf(constants.ErrDeleteData)
		}
	}

	err = objectGroupRepo.DeleteMany(context.Background(), bson.M{"group": group})
	if err != nil {
		fmt.Errorf(constants.ErrDeleteData)
	}

	studentRepo := mongodb.NewRepository[models.StudentsModel, struct{}](db.Collection(constants.StudentCollection))
	err = studentRepo.DeleteMany(context.Background(), bson.M{"group": group})
	if err != nil {
		fmt.Errorf(constants.ErrDeleteData)
	}

	groupRepo := mongodb.NewRepository[models.GroupsModel, struct{}](db.Collection(constants.GroupCollection))
	err = groupRepo.DeleteOne(context.Background(), bson.M{"_id": group})
	if err != nil {
		fmt.Errorf(constants.ErrDeleteData)
	}

	return nil
}

func UpGroup(c *fiber.Ctx, db *mongo.Database) error {
	groupRepo := mongodb.NewRepository[models.GroupsModel, struct{}](db.Collection(constants.GroupCollection))
	groupFindAll, err := groupRepo.FindAll(context.Background(), bson.M{})
	if err != nil {
		services.Logging(db, "/api/admin/up_group", "PATCH", "400", nil, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrGroupNotFound,
		})
	}

	var errMsg []string

	for _, group := range groupFindAll {
		var tempGroup string
		var numDash string
		flagDash := false

		if strings.Contains(group.Group, "-") {
			splitGroup := strings.Split(group.Group, "-")
			numDash = splitGroup[0]
			tempGroup = splitGroup[1]
			flagDash = true
		} else {
			tempGroup = group.Group
		}

		intGroup, err := strconv.Atoi(tempGroup)
		if err != nil {
			errMsg = append(errMsg, err.Error())
			continue
		}

		if CheckLastGroup(intGroup) {
			DeleteUserData(db, group.ID)
		}

		intGroup += 10

		if flagDash {
			tempGroup = fmt.Sprintf("%v-%v", numDash, intGroup)
		} else {
			tempGroup = fmt.Sprintf("%v", intGroup)
		}

		err = groupRepo.UpdateOne(context.Background(), bson.M{"_id": group.ID}, bson.M{"$set": bson.M{"group": tempGroup}})
		if err != nil {
			errMsg = append(errMsg, err.Error())
		}
	}

	if len(errMsg) != 0 {
		services.Logging(db, "/api/admin/up_group", "PATCH", "400", nil, errMsg)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errMsg,
		})
	}

	services.Logging(db, "/api/admin/up_group", "PATCH", "200", nil, constants.SuccGroupUp)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccGroupUp,
	})
}

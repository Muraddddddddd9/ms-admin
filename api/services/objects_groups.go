package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ms-admin/api/constants"
	"ms-admin/api/core"
	"ms-admin/api/utils"
	"strings"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateObjectsGroups(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var objectsGroups models.ObjectsGroupsModel
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&objectsGroups); err != nil {
		return nil, fmt.Errorf("%v: %v", constants.ErrInvalidDataObjectForGroup, err)
	}

	checkReferences := []core.ReferenceCheck{
		{Collection: constants.ObjectCollection, ID: objectsGroups.Object, ErrMsg: constants.ErrObjectNotFound},
		{Collection: constants.GroupCollection, ID: objectsGroups.Group, ErrMsg: constants.ErrGroupNotFound},
		{Collection: constants.TeacherCollection, ID: objectsGroups.Teacher, ErrMsg: constants.ErrTeacherNotFound},
	}

	return core.CreateDocument[*core.ObjectsGroupsModel](
		db,
		data,
		constants.ObjectGroupCollection,
		bson.M{"object": objectsGroups.Object, "group": objectsGroups.Group},
		checkReferences,
	)
}

func ReadObjectsGroups(db *mongo.Database, page, pageSize int) (map[string]interface{}, error) {
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         constants.ObjectCollection,
				"localField":   "object",
				"foreignField": "_id",
				"as":           "objectsData",
			},
		},
		{
			"$lookup": bson.M{
				"from":         constants.GroupCollection,
				"localField":   "group",
				"foreignField": "_id",
				"as":           "groupsData",
			},
		},
		{
			"$lookup": bson.M{
				"from":         constants.TeacherCollection,
				"localField":   "teacher",
				"foreignField": "_id",
				"as":           "teachersData",
			},
		},
		{
			"$unwind": "$objectsData",
		},
		{
			"$unwind": "$groupsData",
		},
		{
			"$unwind": "$teachersData",
		},
		{
			"$project": bson.M{
				"_id":    1,
				"object": "$objectsData.object",
				"group":  "$groupsData.group",
				"teacher": bson.M{
					"$concat": bson.A{
						"$teachersData.name",
						" ",
						"$teachersData.surname",
						" ",
						"$teachersData.patronymic",
					},
				},
			},
		},
		{
			"$sort": bson.M{
				"group": 1,
			},
		},
	}

	return core.ReadAggregateDocument[models.ObjectsGroupsWithGroupAndTeacherModel](
		db,
		constants.ObjectGroupCollection,
		pipeline,
		[]string{"ID"},
		[]string{constants.TeacherCollection, constants.ObjectCollection, constants.GroupCollection},
		page, pageSize,
	)
}

func GetAllObject(c *fiber.Ctx, db *mongo.Database) error {
	group := strings.TrimSpace(c.Params("group"))
	groupID, err := primitive.ObjectIDFromHex(group)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constants.ErrServerError,
		})
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"group": groupID,
			},
		},
		{
			"$lookup": bson.M{
				"from":         constants.ObjectCollection,
				"localField":   "object",
				"foreignField": "_id",
				"as":           "objectData",
			},
		},
		{
			"$unwind": "$objectData",
		},
		{
			"$project": bson.M{
				"_id":     1,
				"object":  "$objectData.object",
				"group":   1,
				"teacher": 1,
			},
		},
	}

	objectsGroupRepo := mongodb.NewRepository[struct{}, models.ObjectsGroupsWithGroupAndTeacherModel](db.Collection(constants.ObjectGroupCollection))
	objectsAggregateAll, err := objectsGroupRepo.AggregateAll(context.Background(), pipeline)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": constants.ErrObjectNotFound,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"objects": objectsAggregateAll,
	})
}

func DeleteObjectsGroups(db *mongo.Database, collectionName string, ids []primitive.ObjectID) (string, error) {
	checkReferencesOther := []core.ReferenceCheckOther{
		{Collection: constants.EvaluationCollection, Field: "object"},
	}

	return core.DeleteDocument[models.ObjectsGroupsModel](
		db,
		collectionName,
		ids,
		checkReferencesOther,
	)
}

func UpdateObjectsGroups(
	db *mongo.Database,
	rdb *redis.Client,
	collection string,
	id primitive.ObjectID,
	label string,
	newData string,
) error {
	objectGroupsRepo := mongodb.NewRepository[models.ObjectsGroupsModel, struct{}](db.Collection(constants.ObjectGroupCollection))
	objectGroupsFindOne, err := objectGroupsRepo.FindOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf(constants.ErrObjectNotFound)
	}

	if label == "object" {
		newDataObject, _ := primitive.ObjectIDFromHex(newData)
		err := utils.CheckReplica(db, constants.ObjectGroupCollection, bson.M{"group": objectGroupsFindOne.Group, label: newDataObject})
		if err != nil {
			return err
		}
	} else if label == "group" {
		newDataObject, _ := primitive.ObjectIDFromHex(newData)
		err := utils.CheckReplica(db, constants.ObjectGroupCollection, bson.M{"object": objectGroupsFindOne.Object, label: newDataObject})
		if err != nil {
			return err
		}
	}

	return core.UpdateDocument[models.ObjectsGroupsModel, struct{}](
		db,
		rdb,
		id,
		constants.ObjectGroupCollection,
		label,
		newData,
		[]string{"object", "group", "teacher"},
		nil,
	)
}

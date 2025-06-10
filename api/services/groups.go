package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ms-admin/api/constants"
	"ms-admin/api/core"
	"ms-admin/api/utils"
	"strings"

	"github.com/Muraddddddddd9/ms-database/models"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateGroups(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var group models.GroupsModel
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&group); err != nil {
		return nil, fmt.Errorf("%v: %v", constants.ErrInvalidDataGroup, err)
	}

	checkReferences := []core.ReferenceCheck{
		{Collection: constants.TeacherCollection, ID: group.Teacher, ErrMsg: constants.ErrTeacherNotFound},
	}

	return core.CreateDocument[*core.GroupsModel](
		db,
		data,
		constants.GroupCollection,
		bson.M{"group": strings.TrimSpace(strings.ToLower(group.Group))},
		checkReferences,
	)
}

func ReadGroups(db *mongo.Database, page, pageSize int) (map[string]interface{}, error) {
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         constants.TeacherCollection,
				"localField":   "teacher",
				"foreignField": "_id",
				"as":           "teacherData",
			},
		},
		{
			"$unwind": "$teacherData",
		},
		{
			"$project": bson.M{
				"_id":   1,
				"group": 1,
				"teacher": bson.M{
					"$concat": bson.A{
						"$teacherData.name",
						" ",
						"$teacherData.surname",
						" ",
						"$teacherData.patronymic",
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

	return core.ReadAggregateDocument[models.GroupsWithTeacherModel](
		db,
		constants.GroupCollection,
		pipeline,
		[]string{"ID"},
		[]string{constants.TeacherCollection},
		page, pageSize,
	)
}

func DeleteGroups(db *mongo.Database, collectionName string, ids []primitive.ObjectID) (string, error) {
	checkReferencesOther := []core.ReferenceCheckOther{
		{Collection: constants.ObjectGroupCollection, Field: "group"},
		{Collection: constants.StudentCollection, Field: "group"},
	}

	return core.DeleteDocument[models.GroupsModel](
		db,
		collectionName,
		ids,
		checkReferencesOther,
	)
}

func UpdateGroups(
	db *mongo.Database,
	rdb *redis.Client,
	collection string,
	id primitive.ObjectID,
	label string,
	newData string,
) error {
	if label == "group" {
		err := utils.CheckReplica(db, constants.GroupCollection, bson.M{label: newData})
		if err != nil {
			return err
		}
	}

	return core.UpdateDocument[models.GroupsModel, struct{}](
		db,
		rdb,
		id,
		constants.GroupCollection,
		label,
		newData,
		[]string{"teacher"},
		nil,
	)
}

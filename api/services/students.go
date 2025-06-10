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
	"golang.org/x/crypto/bcrypt"
)

func CreateStudents(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var student models.StudentsModel
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&student); err != nil {
		return nil, fmt.Errorf("%v: %v", constants.ErrInvalidDataStudent, err)
	}

	checkReferences := []core.ReferenceCheck{
		{Collection: constants.GroupCollection, ID: student.Group, ErrMsg: constants.ErrGroupNotFound},
		{Collection: constants.StatusCollection, ID: student.Status, ErrMsg: constants.ErrStatusNotFound},
	}

	return core.CreateDocument[*core.StudentsModel](
		db,
		data,
		constants.StudentCollection,
		bson.M{"email": strings.TrimSpace(strings.ToLower(student.Email))},
		checkReferences,
	)
}

func ReadStudents(db *mongo.Database, page, pageSize int) (map[string]interface{}, error) {
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         constants.GroupCollection,
				"localField":   "group",
				"foreignField": "_id",
				"as":           "groupData",
			},
		},
		{
			"$lookup": bson.M{
				"from":         constants.StatusCollection,
				"localField":   "status",
				"foreignField": "_id",
				"as":           "statusData",
			},
		},
		{
			"$unwind": "$groupData",
		},
		{
			"$unwind": "$statusData",
		},
		{
			"$project": bson.M{
				"name":       1,
				"surname":    1,
				"patronymic": 1,
				"group":      "$groupData.group",
				"email":      1,
				"password":   1,
				"telegram":   1,
				"diplomas":   1,
				"ips":        1,
				"status":     "$statusData.status",
			},
		},
		{
			"$sort": bson.M{
				"group": 1,
			},
		},
	}

	return core.ReadAggregateDocument[models.StudentsWithGroupAndStatusModel](
		db,
		constants.StudentCollection,
		pipeline,
		[]string{"ID", "Telegram", "Diplomas"},
		[]string{constants.StatusCollection, constants.GroupCollection},
		page, pageSize,
	)
}

func DeleteStudents(db *mongo.Database, collectionName string, ids []primitive.ObjectID) (string, error) {
	checkReferencesOther := []core.ReferenceCheckOther{
		{Collection: constants.EvaluationCollection, Field: "student"},
	}

	return core.DeleteDocument[models.ObjectsModel](
		db,
		collectionName,
		ids,
		checkReferencesOther,
	)
}

func UpdateStudents(
	db *mongo.Database,
	rdb *redis.Client,
	collection string,
	id primitive.ObjectID,
	label string,
	newData string,
) error {
	if label == "email" || label == "telegram" {
		newData = strings.ToLower(newData)
		err := utils.CheckReplica(db, constants.StudentCollection, bson.M{label: newData})
		if err != nil {
			return err
		}
	} else if label == "password" {
		newData = strings.TrimSpace(newData)
		bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(newData), bcrypt.DefaultCost)
		newData = string(bcryptPassword)
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": id,
			},
		},
		{
			"$lookup": bson.M{
				"from":         constants.GroupCollection,
				"localField":   "group",
				"foreignField": "_id",
				"as":           "groupData",
			},
		},
		{
			"$lookup": bson.M{
				"from":         constants.StatusCollection,
				"localField":   "status",
				"foreignField": "_id",
				"as":           "statusData",
			},
		},
		{
			"$unwind": "$groupData",
		},
		{
			"$unwind": "$statusData",
		},
		{
			"$project": bson.M{
				"name":       1,
				"surname":    1,
				"patronymic": 1,
				"group":      "$groupData.group",
				"email":      1,
				"telegram":   1,
				"diplomas":   1,
				"ips":        1,
				"status":     "$statusData.status",
			},
		},
	}

	return core.UpdateDocument[models.StudentsModel, models.StudentsWithGroupAndStatusModel](
		db,
		rdb,
		id,
		constants.StudentCollection,
		label,
		newData,
		[]string{"group", "status"},
		pipeline,
	)
}

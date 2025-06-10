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

func CreateTeachers(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var teacher models.TeachersModel
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&teacher); err != nil {
		return nil, fmt.Errorf("%v: %v", constants.ErrInvalidDataTeacher, err)
	}

	checkReferences := []core.ReferenceCheck{
		{Collection: constants.StatusCollection, ID: teacher.Status, ErrMsg: constants.ErrStatusNotFound},
	}

	return core.CreateDocument[*core.TeachersModel](
		db,
		data,
		constants.TeacherCollection,
		bson.M{"email": strings.TrimSpace(strings.ToLower(teacher.Email))},
		checkReferences,
	)
}

func ReadTeachers(db *mongo.Database, page, pageSize int) (map[string]interface{}, error) {
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         constants.StatusCollection,
				"localField":   "status",
				"foreignField": "_id",
				"as":           "statusData",
			},
		},
		{
			"$unwind": "$statusData",
		},
		{
			"$project": bson.M{
				"name":       1,
				"surname":    1,
				"patronymic": 1,
				"email":      1,
				"password":   1,
				"telegram":   1,
				"ips":        1,
				"status":     "$statusData.status",
			},
		},
		{
			"$sort": bson.M{
				"name": 1,
			},
		},
	}

	return core.ReadAggregateDocument[models.TeachersWithStatusModel](
		db,
		constants.TeacherCollection,
		pipeline,
		[]string{"ID", "Telegram", "Diplomas"},
		[]string{constants.StatusCollection},
		page, pageSize,
	)
}

func DeleteTeachers(db *mongo.Database, collectionName string, ids []primitive.ObjectID) (string, error) {
	checkReferencesOther := []core.ReferenceCheckOther{
		{Collection: constants.ObjectGroupCollection, Field: "teacher"},
		{Collection: constants.GroupCollection, Field: "teacher"},
		{Collection: constants.ContestCollection, Field: "teacher"},
	}

	return core.DeleteDocument[models.ObjectsModel](
		db,
		collectionName,
		ids,
		checkReferencesOther,
	)
}

func UpdateTeachers(
	db *mongo.Database,
	rdb *redis.Client,
	collection string,
	id primitive.ObjectID,
	label string,
	newData string,
) error {
	if label == "email" || label == "telegram" {
		newData = strings.ToLower(newData)
		err := utils.CheckReplica(db, constants.TeacherCollection, bson.M{label: newData})
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
			"$lookup": bson.M{
				"from":         constants.StatusCollection,
				"localField":   "status",
				"foreignField": "_id",
				"as":           "statusData",
			},
		},
		{
			"$unwind": "$statusData",
		},
		{
			"$project": bson.M{
				"name":       1,
				"surname":    1,
				"patronymic": 1,
				"email":      1,
				"telegram":   1,
				"ips":        1,
				"status":     "$statusData.status",
			},
		},
	}

	return core.UpdateDocument[models.TeachersModel, models.TeachersWithStatusModel](
		db,
		rdb,
		id,
		constants.TeacherCollection,
		label,
		newData,
		[]string{"status"},
		pipeline,
	)
}

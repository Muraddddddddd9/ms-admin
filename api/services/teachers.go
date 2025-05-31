package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ms-admin/api/constants"
	"ms-admin/api/core"
	"strings"

	"github.com/Muraddddddddd9/ms-database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func ReadTeachers(db *mongo.Database) (map[string]interface{}, error) {
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
	}

	return core.ReadAggregateDocument[models.TeachersWithStatusModel](
		db,
		constants.TeacherCollection,
		pipeline,
		[]string{"ID", "Telegram", "Diplomas", "IPs"},
		[]string{constants.StatusCollection},
	)
}

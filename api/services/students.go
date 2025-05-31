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

func ReadStudents(db *mongo.Database) (map[string]interface{}, error) {
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
		[]string{"ID", "Telegram", "Diplomas", "IPs"},
		[]string{constants.StatusCollection, constants.GroupCollection},
	)
}

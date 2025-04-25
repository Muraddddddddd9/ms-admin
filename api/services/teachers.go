package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ms-admin/api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	TeacherCollection = "teachers"
)

func CreateTeachers(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var teacher models.TeachersModel
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&teacher); err != nil {
		return nil, fmt.Errorf("%v: %v", "Неверные данные учителя", err)
	}

	err := CheckReplica(db, TeacherCollection, bson.M{"email": teacher.Email})
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return teacher, nil
}

func ReadTeachers(db *mongo.Database) (interface{}, []string, interface{}, error) {
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "statuses",
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

	cursor, err := db.Collection(TeacherCollection).Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}
	defer cursor.Close(context.TODO())

	var results []models.TeachersWithStatusModel
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}

	var structForHead models.TeachersWithStatusModel
	header := GetFieldNames(structForHead)
	header = FilterHeaders(header, []string{"ID", "Telegram", "Diplomas", "IPs"})

	var selectResult models.SelectModels
	err = SelectData(db, StatusCollection, &selectResult)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}

	return results, header, selectResult, nil
}

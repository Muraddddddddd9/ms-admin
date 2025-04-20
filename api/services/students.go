package services

import (
	"context"
	"encoding/json"
	"fmt"
	"ms-admin/api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	StidentCollection = "students"
)

func CreateStudent(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var student models.StudentsModel
	if err := json.Unmarshal(data, &student); err != nil {
		return nil, fmt.Errorf("%s", "Неверные данные студента")
	}

	err := CheckReplica(db, StidentCollection, bson.M{"email": student.Email})
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	var findGroup any
	err = db.Collection("groups").FindOne(context.TODO(), bson.M{"_id": student.Group}).Decode(&findGroup)
	if err != nil {
		return nil, fmt.Errorf("%s", "Группа не найдена")
	}

	return student, nil
}

func ReadStudent(db *mongo.Database) (interface{}, error) {
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "groups",
				"localField":   "group",
				"foreignField": "_id",
				"as":           "groupData",
			},
		},
		{
			"$unwind": "$groupData",
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
				"status":     1,
			},
		},
	}

	cursor, err := db.Collection(StidentCollection).Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer cursor.Close(context.TODO())

	var results []models.StudentWithGroup
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return results, nil
}

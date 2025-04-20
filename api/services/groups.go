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
	GroupCollection = "groups"
)

func CreateGroups(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var group models.GroupsModel
	if err := json.Unmarshal(data, &group); err != nil {
		return nil, fmt.Errorf("%s", "Неверные данные группы")
	}

	err := CheckReplica(db, GroupCollection, bson.M{"group": group.Group})
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	var findTeacher any
	err = db.Collection(TeacherCollection).FindOne(context.TODO(), bson.M{"_id": group.TeacherId}).Decode(&findTeacher)
	if err != nil {
		return nil, fmt.Errorf("%s", "Учитель не найден")
	}

	return group, nil
}

func ReadGroups(db *mongo.Database) (interface{}, error) {
	pipline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "teachers",
				"localField":   "teacher_id",
				"foreignField": "_id",
				"as":           "teacherData",
			},
		},
		{
			"$unwind": "$teacherData",
		},
		{
			"$project": bson.M{
				"_id":        1,
				"group":      1,
				"teacher_id": bson.M{
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
	}

	cursor, err := db.Collection(GroupCollection).Aggregate(context.TODO(), pipline)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer cursor.Close(context.TODO())

	var results []models.GroupWithTeacher
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return results, nil
}

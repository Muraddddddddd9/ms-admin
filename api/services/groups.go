package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ms-admin/api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	GroupCollection = "groups"
)

func CreateGroups(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var group models.GroupsModel
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&group); err != nil {
		return nil, fmt.Errorf("%v: %v", "Неверные данные группы", err)
	}

	if group.Group == "" {
		return nil, fmt.Errorf("поле 'group' не может быть пустым")
	}
	
	if group.Teacher == primitive.NilObjectID {
		return nil, fmt.Errorf("поле 'teacher' не может быть пустым")
	}

	err := CheckReplica(db, GroupCollection, bson.M{"group": group.Group})
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	var findTeacher any
	err = db.Collection(TeacherCollection).FindOne(context.TODO(), bson.M{"_id": group.Teacher}).Decode(&findTeacher)
	if err != nil {
		return nil, fmt.Errorf("%s", "Учитель не найден")
	}

	return group, nil
}

func ReadGroups(db *mongo.Database) (interface{}, []string, interface{}, error) {
	pipline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "teachers",
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
	}

	cursor, err := db.Collection(GroupCollection).Aggregate(context.TODO(), pipline)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}
	defer cursor.Close(context.TODO())

	var results []models.GroupsWithTeacherModel
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}

	var structForHead models.GroupsWithTeacherModel
	header := GetFieldNames(structForHead)
	header = FilterHeaders(header, []string{"ID"})

	var selectResult models.SelectModels
	err = SelectData(db, TeacherCollection, &selectResult)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}

	return results, header, selectResult, nil
}

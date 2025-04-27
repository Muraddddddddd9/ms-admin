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
	ObjectGroupCollection = "objects_groups"
)

func CreateObjectsGroups(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var objectsGroups models.ObjectsGroupsModel
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&objectsGroups); err != nil {
		return nil, fmt.Errorf("%v: %v", "Неверные данные предмета для группы", err)
	}

	err := CheckReplica(db, ObjectGroupCollection, bson.M{"object": objectsGroups.Object, "group": objectsGroups.Group})
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	var findObject any
	err = db.Collection(ObjectCollection).FindOne(context.TODO(), bson.M{"_id": objectsGroups.Object}).Decode(&findObject)
	if err != nil {
		return nil, fmt.Errorf("%s", "Предмет не найдена")
	}

	var findGroup any
	err = db.Collection(GroupCollection).FindOne(context.TODO(), bson.M{"_id": objectsGroups.Group}).Decode(&findGroup)
	if err != nil {
		return nil, fmt.Errorf("%s", "Группа не найдена")
	}

	var findTeacher any
	err = db.Collection(TeacherCollection).FindOne(context.TODO(), bson.M{"_id": objectsGroups.Teacher}).Decode(&findTeacher)
	if err != nil {
		return nil, fmt.Errorf("%s", "Учитель не найдена")
	}

	return objectsGroups, nil
}

func ReadObjectsGroups(db *mongo.Database) (interface{}, []string, interface{}, error) {
	pipline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "objects",
				"localField":   "object",
				"foreignField": "_id",
				"as":           "objectsData",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "groups",
				"localField":   "group",
				"foreignField": "_id",
				"as":           "groupsData",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "teachers",
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
	}

	cursor, err := db.Collection(ObjectGroupCollection).Aggregate(context.TODO(), pipline)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}
	defer cursor.Close(context.TODO())

	var results []models.ObjectsGroupsWithGroupAndTeacherModel
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}

	var structForHead models.ObjectsGroupsWithGroupAndTeacherModel
	header := GetFieldNames(structForHead)
	header = FilterHeaders(header, []string{"ID"})

	var selectResult models.SelectModels
	err = SelectData(db, ObjectCollection, &selectResult)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}

	err = SelectData(db, GroupCollection, &selectResult)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}

	err = SelectData(db, TeacherCollection, &selectResult)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}

	return results, header, selectResult, nil
}

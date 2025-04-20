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
	ObjectGroupCollection = "objects_groups"
)

func CreateObjectsGroups(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var objectsGroups models.ObjectsGroupsModel
	if err := json.Unmarshal(data, &objectsGroups); err != nil {
		return nil, fmt.Errorf("%s", "Неверные данные предмета для группы")
	}

	err := CheckReplica(db, ObjectGroupCollection, bson.M{"object": objectsGroups.ObjectId})
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	var findObject any
	err = db.Collection(ObjectCollection).FindOne(context.TODO(), bson.M{"_id": objectsGroups.ObjectId}).Decode(&findObject)
	if err != nil {
		return nil, fmt.Errorf("%s", "Предмет не найдена")
	}

	var findGroup any
	err = db.Collection(GroupCollection).FindOne(context.TODO(), bson.M{"_id": objectsGroups.GroupId}).Decode(&findGroup)
	if err != nil {
		return nil, fmt.Errorf("%s", "Группа не найдена")
	}

	var findTeacher any
	err = db.Collection(TeacherCollection).FindOne(context.TODO(), bson.M{"_id": objectsGroups.TeacherId}).Decode(&findTeacher)
	if err != nil {
		return nil, fmt.Errorf("%s", "Учитель не найдена")
	}

	return objectsGroups, nil
}

func ReadObjectsGroups(db *mongo.Database) (interface{}, error) {
	pipline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "objects",
				"localField":   "object_id",
				"foreignField": "_id",
				"as":           "objectsData",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "groups",
				"localField":   "group_id",
				"foreignField": "_id",
				"as":           "groupsData",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "teachers",
				"localField":   "teacher_id",
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
				"_id":       1,
				"object_id": "$objectsData.object",
				"group_id":  "$groupsData.group",
				"teacher_id": bson.M{
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
		return nil, fmt.Errorf("%v", err)
	}
	defer cursor.Close(context.TODO())

	var results []models.ObjectsGroupsWithGroupAndTeacherModel
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return results, nil
}

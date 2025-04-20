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
	TeacherCollection = "teachers"
)

func CreateTeachers(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var teacher models.TeachersModel
	if err := json.Unmarshal(data, &teacher); err != nil {
		return nil, fmt.Errorf("%s", "Неверные данные учителя")
	}

	err := CheckReplica(db, TeacherCollection, bson.M{"email": teacher.Email})
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return teacher, nil
}

func ReadTeachers(db *mongo.Database) (interface{}, error) {
	var results []models.TeachersModel
	cursor, err := db.Collection(TeacherCollection).Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return results, nil
}

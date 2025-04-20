package services

import (
	"context"
	"encoding/json"
	"fmt"
	"ms-admin/api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateStudent(db *mongo.Database, data Data) (interface{}, error) {
	var student models.StudentsModel
	if err := json.Unmarshal(data.NewData, &student); err != nil {
		return nil, fmt.Errorf("%s", "Неверные данные студента")
	}

	err := checkReplica(db, data, bson.M{"email": student.Email})
	if err != nil {
		return nil, fmt.Errorf("%e", err)
	}

	student.Group, err = primitive.ObjectIDFromHex(student.Group.String())
	if err != nil {
		return nil, fmt.Errorf("%e", err)
	}

	var findGroup any
	err = db.Collection("groups").FindOne(context.TODO(), bson.M{"_id": student.Group}).Decode(&findGroup)
	if err != nil {
		return nil, fmt.Errorf("%s", "Группа не найдена")
	}

	return student, nil
}

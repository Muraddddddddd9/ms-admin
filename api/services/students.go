package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ms-admin/api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	StidentCollection = "students"
)

func CreateStudents(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var student models.StudentsModel
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&student); err != nil {
		return nil, fmt.Errorf("%v: %v", "Неверные данные студента", err)
	}

	if student.Name == "" {
		return nil, fmt.Errorf("поле 'name' не может быть пустым")
	}

	if student.Surname == "" {
		return nil, fmt.Errorf("поле 'surname' не может быть пустым")
	}

	if student.Patronymic == "" {
		return nil, fmt.Errorf("поле 'patronymic' не может быть пустым")
	}

	if student.Email == "" {
		return nil, fmt.Errorf("поле 'email' не может быть пустым")
	}

	if student.Password == "" {
		return nil, fmt.Errorf("поле 'password' не может быть пустым")
	}

	err := CheckReplica(db, StidentCollection, bson.M{"email": student.Email})
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	student.Diplomas = []string{}
	student.IPs = []string{}

	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(student.Password), bcrypt.DefaultCost)
	student.Password = string(bcryptPassword)

	var findGroup any
	err = db.Collection(GroupCollection).FindOne(context.TODO(), bson.M{"_id": student.Group}).Decode(&findGroup)
	if err != nil {
		return nil, fmt.Errorf("%s", "Группа не найдена")
	}

	var findStatus any
	err = db.Collection(StatusCollection).FindOne(context.TODO(), bson.M{"_id": student.Status}).Decode(&findStatus)
	if err != nil {
		return nil, fmt.Errorf("%s", "Статус не найдена")
	}

	return student, nil
}

func ReadStudents(db *mongo.Database) (interface{}, []string, interface{}, error) {
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
			"$lookup": bson.M{
				"from":         "statuses",
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
	}

	cursor, err := db.Collection(StidentCollection).Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}
	defer cursor.Close(context.TODO())

	var results []models.StudentsWithGroupAndStatusModel
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}

	var structForHead models.StudentsWithGroupAndStatusModel
	header := GetFieldNames(structForHead)
	header = FilterHeaders(header, []string{"ID", "Telegram", "Diplomas", "IPs"})

	var selectResult models.SelectModels
	err = SelectData(db, GroupCollection, &selectResult)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}

	err = SelectData(db, StatusCollection, &selectResult)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}

	return results, header, selectResult, nil
}

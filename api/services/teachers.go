package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ms-admin/api/messages"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	TeacherCollection = "teachers"
)

func CreateTeachers(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var teacher models.TeachersModel

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&teacher); err != nil {
		return nil, fmt.Errorf("%v: %v", messages.ErrInvalidDataTeacher, err)
	}

	fields := map[string]string{
		"name":       teacher.Name,
		"surname":    teacher.Surname,
		"patronymic": teacher.Patronymic,
		"email":      teacher.Email,
		"password":   teacher.Password,
	}

	for name, value := range fields {
		if value == "" {
			return nil, fmt.Errorf(messages.ErrFieldCannotEmpty, name)
		}
	}

	err := CheckReplica(db, TeacherCollection, bson.M{"email": teacher.Email})
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	teacher.IPs = []string{}

	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(teacher.Password), bcrypt.DefaultCost)
	teacher.Password = string(bcryptPassword)

	statusRepo := mongodb.NewRepository[models.StatusesModel, interface{}](db.Collection(StatusCollection))
	_, err = statusRepo.FindOne(context.Background(), bson.M{"_id": teacher.Status})
	if err != nil {
		return nil, fmt.Errorf("%s", messages.ErrStatusNotFound)
	}

	teacherRepo := mongodb.NewRepository[models.TeachersModel, models.TeachersWithStatusModel](db.Collection(TeacherCollection))
	teacherId, err := teacherRepo.InsertOne(context.Background(), &teacher)
	if err != nil {
		return nil, err
	}

	return teacherId, nil
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

	teacherRepo := mongodb.NewRepository[models.TeachersModel, models.TeachersWithStatusModel](db.Collection(TeacherCollection))
	teacherAggregate, err := teacherRepo.AggregateAll(context.Background(), pipeline)
	if err != nil {
		return nil, nil, nil, err
	}

	var structForHead models.TeachersWithStatusModel
	header := GetFieldNames(structForHead)
	header = FilterHeaders(header, []string{"ID", "Telegram", "Diplomas", "IPs"})

	var selectResult models.SelectModels
	err = SelectData(db, StatusCollection, &selectResult)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}

	return teacherAggregate, header, selectResult, nil
}

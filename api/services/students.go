package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ms-admin/api/constants"
	"strings"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func CreateStudents(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var student models.StudentsModel

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&student); err != nil {
		return nil, fmt.Errorf("%v: %v", constants.ErrInvalidDataStudent, err)
	}

	student.Name = strings.TrimSpace(student.Name)
	student.Surname = strings.TrimSpace(student.Surname)
	student.Patronymic = strings.TrimSpace(student.Patronymic)
	student.Email = strings.TrimSpace(strings.ToLower(student.Email))
	student.Password = strings.TrimSpace(student.Password)

	fields := map[string]string{
		"name":       student.Name,
		"surname":    student.Surname,
		"patronymic": student.Patronymic,
		"email":      student.Email,
		"password":   student.Password,
	}

	for name, value := range fields {
		if strings.TrimSpace(value) == "" {
			return nil, fmt.Errorf(constants.ErrFieldCannotEmpty, name)
		}
	}

	err := CheckReplica(db, constants.StudentCollection, bson.M{"email": student.Email})
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	student.Diplomas = []string{}
	student.IPs = []string{}

	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(student.Password), bcrypt.DefaultCost)
	student.Password = string(bcryptPassword)

	groupRepo := mongodb.NewRepository[models.GroupsModel, struct{}](db.Collection(constants.GroupCollection))
	_, err = groupRepo.FindOne(context.Background(), bson.M{"_id": student.Group})
	if err != nil {
		return nil, fmt.Errorf("%s", constants.ErrGroupNotFound)
	}

	statusRepo := mongodb.NewRepository[models.StatusesModel, struct{}](db.Collection(constants.StatusCollection))
	_, err = statusRepo.FindOne(context.Background(), bson.M{"_id": student.Status})
	if err != nil {
		return nil, fmt.Errorf("%s", constants.ErrStatusNotFound)
	}

	studentRepo := mongodb.NewRepository[models.StudentsModel, struct{}](db.Collection(constants.StudentCollection))
	studentID, err := studentRepo.InsertOne(context.Background(), &student)
	if err != nil {
		return nil, err
	}

	return studentID, nil
}

func ReadStudents(db *mongo.Database) (interface{}, []string, interface{}, error) {
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         constants.GroupCollection,
				"localField":   "group",
				"foreignField": "_id",
				"as":           "groupData",
			},
		},
		{
			"$lookup": bson.M{
				"from":         constants.StatusCollection,
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

	studentRepo := mongodb.NewRepository[struct{}, models.StudentsWithGroupAndStatusModel](db.Collection(constants.StudentCollection))
	studentAggregate, err := studentRepo.AggregateAll(context.Background(), pipeline)
	if err != nil {
		return nil, nil, nil, err
	}

	var structForHead models.StudentsWithGroupAndStatusModel
	header := GetFieldNames(structForHead)
	header = FilterHeaders(header, []string{"ID", "Telegram", "Diplomas", "IPs"})

	var selectResult models.SelectModels
	err = SelectData(db, constants.GroupCollection, &selectResult)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}

	err = SelectData(db, constants.StatusCollection, &selectResult)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}

	return studentAggregate, header, selectResult, nil
}

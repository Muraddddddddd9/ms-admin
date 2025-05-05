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
)

func CreateGroups(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var group models.GroupsModel

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&group); err != nil {
		return nil, fmt.Errorf("%v: %v", constants.ErrInvalidDataGroup, err)
	}

	group.Group = strings.ToLower(strings.TrimSpace(group.Group))

	fileds := map[string]string{
		"group": group.Group,
	}

	for name, value := range fileds {
		if value == "" {
			return nil, fmt.Errorf(constants.ErrFieldCannotEmpty, name)
		}
	}

	err := CheckReplica(db, constants.GroupCollection, bson.M{"group": group.Group})
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	teacherRepo := mongodb.NewRepository[models.TeachersModel, models.TeachersWithStatusModel](db.Collection(constants.TeacherCollection))
	_, err = teacherRepo.FindOne(context.Background(), bson.M{"_id": group.Teacher})
	if err != nil {
		return nil, fmt.Errorf("%s", constants.ErrTeacherNotFound)
	}

	groupRepo := mongodb.NewRepository[models.GroupsModel, models.GroupsWithTeacherModel](db.Collection(constants.GroupCollection))
	groupID, err := groupRepo.InsertOne(context.Background(), &group)
	if err != nil {
		return nil, err
	}

	return groupID, nil
}

func ReadGroups(db *mongo.Database) (interface{}, []string, interface{}, error) {
	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         constants.TeacherCollection,
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

	groupRepo := mongodb.NewRepository[models.GroupsModel, models.GroupsWithTeacherModel](db.Collection(constants.GroupCollection))
	groupAggregate, err := groupRepo.AggregateAll(context.Background(), pipeline)
	if err != nil {
		return nil, nil, nil, err
	}

	var structForHead models.GroupsWithTeacherModel
	header := GetFieldNames(structForHead)
	header = FilterHeaders(header, []string{"ID"})

	var selectResult models.SelectModels
	err = SelectData(db, constants.TeacherCollection, &selectResult)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%v", err)
	}

	return groupAggregate, header, selectResult, nil
}

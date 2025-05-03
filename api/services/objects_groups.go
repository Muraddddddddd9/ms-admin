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
)

var (
	ObjectGroupCollection = "objects_groups"
)

func CreateObjectsGroups(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var objectsGroups models.ObjectsGroupsModel

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&objectsGroups); err != nil {
		return nil, fmt.Errorf("%v: %v", messages.ErrInvalidDataObjectForGroup, err)
	}

	err := CheckReplica(db, ObjectGroupCollection, bson.M{"object": objectsGroups.Object, "group": objectsGroups.Group})
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	objectRepo := mongodb.NewRepository[models.ObjectsModel, interface{}](db.Collection(ObjectCollection))
	_, err = objectRepo.FindOne(context.Background(), bson.M{"_id": objectsGroups.Object})
	if err != nil {
		return nil, fmt.Errorf("%s", messages.ErrObjectNotFound)
	}

	groupRepo := mongodb.NewRepository[models.GroupsModel, models.GroupsWithTeacherModel](db.Collection(GroupCollection))
	_, err = groupRepo.FindOne(context.Background(), bson.M{"_id": objectsGroups.Group})
	if err != nil {
		return nil, fmt.Errorf("%s", messages.ErrGroupNotFound)
	}

	teacherRepo := mongodb.NewRepository[models.TeachersModel, models.TeachersWithStatusModel](db.Collection(TeacherCollection))
	_, err = teacherRepo.FindOne(context.Background(), bson.M{"_id": objectsGroups.Teacher})
	if err != nil {
		return nil, fmt.Errorf("%s", messages.ErrTeacherNotFound)
	}

	objectGroupRepo := mongodb.NewRepository[models.ObjectsGroupsModel, models.ObjectsGroupsWithGroupAndTeacherModel](db.Collection(ObjectGroupCollection))
	objectGroupID, err := objectGroupRepo.InsertOne(context.Background(), &objectsGroups)
	if err != nil {
		return nil, err
	}

	return objectGroupID, nil
}

func ReadObjectsGroups(db *mongo.Database) (interface{}, []string, interface{}, error) {
	pipeline := []bson.M{
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

	objectGroupRepo := mongodb.NewRepository[models.ObjectsGroupsModel, models.ObjectsGroupsWithGroupAndTeacherModel](db.Collection(ObjectGroupCollection))
	objectGroupAggregate, err := objectGroupRepo.AggregateAll(context.Background(), pipeline)
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

	return objectGroupAggregate, header, selectResult, nil
}

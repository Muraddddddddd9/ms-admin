package utils

import (
	"context"
	"fmt"
	"ms-admin/api/constants"
	"reflect"

	"github.com/Muraddddddddd9/ms-database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CheckReplica(db *mongo.Database, collection string, filter bson.M) error {
	var result any
	err := db.Collection(collection).FindOne(context.Background(), filter).Decode(&result)
	if err == nil {
		return fmt.Errorf("%s", constants.ErrDataAlreadyExists)
	}

	return nil
}

func GetFieldNames(v interface{}) []string {
	var fields []string
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fields
	}

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		fields = append(fields, typ.Field(i).Name)
	}

	return fields
}

func FilterHeaders(headers []string, toRemove []string) []string {
	removeMap := make(map[string]bool)
	for _, r := range toRemove {
		removeMap[r] = true
	}

	var filtered []string
	for _, h := range headers {
		if !removeMap[h] {
			filtered = append(filtered, h)
		}
	}
	return filtered
}

func SelectData(db *mongo.Database, collection string, selectResult *models.SelectModels) error {
	var res interface{}
	var optionsFind *options.FindOptions

	switch collection {
	case constants.ObjectCollection:
		res = &selectResult.Objects
		optionsFind = options.Find().SetSort(bson.D{{Key: "object", Value: 1}})
	case constants.GroupCollection:
		res = &selectResult.Groups
		optionsFind = options.Find().SetSort(bson.D{{Key: "group", Value: 1}})
	case constants.TeacherCollection:
		res = &selectResult.Teachers
		optionsFind = options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	case constants.StatusCollection:
		res = &selectResult.Statuses
		optionsFind = options.Find().SetSort(bson.D{{Key: "status", Value: 1}})
	}

	cursor, err := db.Collection(collection).Find(context.Background(), bson.M{}, optionsFind)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), res)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func CheckDataOtherTable(db *mongo.Database, collection string, filter bson.M) error {
	err := db.Collection(collection).FindOne(context.Background(), filter).Err()
	if err == nil {
		return fmt.Errorf(constants.ErrDataInCollection, collection)
	}
	if err == mongo.ErrNoDocuments {
		return nil
	}

	return err
}

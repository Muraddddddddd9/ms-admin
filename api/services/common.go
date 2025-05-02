package services

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Muraddddddddd9/ms-database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckReplica(db *mongo.Database, collection string, filter bson.M) error {
	var result any
	err := db.Collection(collection).FindOne(context.TODO(), filter).Decode(&result)
	if err == nil {
		return fmt.Errorf("%s", "Данные уже существуют")
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
	cursor, err := db.Collection(collection).Find(context.TODO(), bson.M{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer cursor.Close(context.TODO())

	var res interface{}
	switch collection {
	case "objects":
		res = &selectResult.Objects
	case "groups":
		res = &selectResult.Groups
	case "teachers":
		res = &selectResult.Teachers
	case "statuses":
		res = &selectResult.Statuses
	}

	err = cursor.All(context.TODO(), res)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func CheckDataOtherTable(db *mongo.Database, collection string, filter bson.M) error {
	err := db.Collection(collection).FindOne(context.TODO(), filter).Err()
	if err == nil {
		return fmt.Errorf("Данные находятся в коллекции %s", collection)
	}
	if err == mongo.ErrNoDocuments {
		return nil
	}

	return err
}

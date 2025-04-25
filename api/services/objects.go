package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ms-admin/api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ObjectCollection = "objects"
)

func CreateObjects(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var object models.ObjectsModel
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&object); err != nil {
		return nil, fmt.Errorf("%v: %v", "Неверные данные предмета", err)
	}

	if object.Object == "" {
		return nil, fmt.Errorf("поле 'object' не может быть пустым")
	}
	
	err := CheckReplica(db, ObjectCollection, bson.M{"object": object.Object})
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return object, nil
}

func ReadObjects(db *mongo.Database) (interface{}, []string, error) {
	cursor, err := db.Collection(ObjectCollection).Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, nil, fmt.Errorf("%v", err)
	}
	defer cursor.Close(context.TODO())

	var results []models.ObjectsModel
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, nil, fmt.Errorf("%v", err)
	}
	var structForHead models.ObjectsModel
	header := GetFieldNames(structForHead)
	header = FilterHeaders(header, []string{"ID"})

	return results, header, nil
}

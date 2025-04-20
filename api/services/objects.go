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
	ObjectCollection = "objects"
)

func CreateObjects(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var object models.ObjectsModel
	if err := json.Unmarshal(data, &object); err != nil {
		return nil, fmt.Errorf("%s", "Неверные данные предмета")
	}

	err := CheckReplica(db, ObjectCollection, bson.M{"object": object.Object})
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return object, nil
}

func ReadObjects(db *mongo.Database) (interface{}, error) {
	cursor, err := db.Collection(ObjectCollection).Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer cursor.Close(context.TODO())

	var results []models.ObjectsModel
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return results, nil
}
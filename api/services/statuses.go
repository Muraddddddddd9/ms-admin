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
	StatusCollection = "statuses"
)

func CreateStatuses(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var status models.StatusesModel
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&status); err != nil {
		return nil, fmt.Errorf("%v: %v", "Неверные данные статуса", err)
	}

	if status.Status == "" {
		return nil, fmt.Errorf("поле 'status' не может быть пустым")
	}

	err := CheckReplica(db, StatusCollection, bson.M{"status": status.Status})
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return status, nil
}

func ReadStatuses(db *mongo.Database) (interface{}, []string, error) {
	cursor, err := db.Collection(StatusCollection).Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, nil, fmt.Errorf("%v", err)
	}
	defer cursor.Close(context.TODO())

	var results []models.StatusesModel
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, nil, fmt.Errorf("%v", err)
	}
	var structForHead models.StatusesModel
	header := GetFieldNames(structForHead)
	header = FilterHeaders(header, []string{"ID"})

	return results, header, nil
}

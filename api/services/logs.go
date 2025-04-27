package services

import (
	"context"
	"fmt"
	"ms-admin/api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	LogsCollection = "log"
)

func ReadLogs(db *mongo.Database) (interface{}, []string, error) {
	cursor, err := db.Collection(LogsCollection).Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, nil, fmt.Errorf("%v", err)
	}
	defer cursor.Close(context.TODO())

	var results []models.Log
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, nil, fmt.Errorf("%v", err)
	}

	var structForHead models.Log
	header := GetFieldNames(structForHead)
	header = FilterHeaders(header, []string{"ID"})

	return results, header, nil
}

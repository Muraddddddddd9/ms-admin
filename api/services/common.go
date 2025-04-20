package services

import (
	"context"
	"fmt"

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

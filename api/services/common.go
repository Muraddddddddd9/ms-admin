package services

import (
	"context"
	"fmt"
	"ms-admin/api/handlers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func checkReplica(db *mongo.Database, data handlers.Data, filter bson.M) error {
	var result any
	err := db.Collection(data.Collection).FindOne(context.TODO(), filter).Decode(&result)
	if err == nil {
		return fmt.Errorf("%s", "Данные уже существуют")
	}

	return nil
}

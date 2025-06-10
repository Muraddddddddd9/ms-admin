package core

import (
	"context"
	"encoding/json"
	"fmt"
	"ms-admin/api/constants"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateDocument[M, T any](
	db *mongo.Database,
	rdb *redis.Client,
	id primitive.ObjectID,
	collectionName string,
	label string,
	newData string,
	arrObject []string,
	pipeline []bson.M,
) error {
	var update bson.M
	if arrObject != nil {
		for _, obj := range arrObject {
			if label == obj {
				newDataObject, err := primitive.ObjectIDFromHex(newData)
				if err != nil {
					return err
				}
				update = bson.M{"$set": bson.M{label: newDataObject}}
				break
			} else {
				update = bson.M{"$set": bson.M{label: newData}}
			}
		}
	} else {
		update = bson.M{"$set": bson.M{label: newData}}
	}

	repo := mongodb.NewRepository[M, struct{}](db.Collection(collectionName))
	err := repo.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf(constants.ErrUpdateData)
	}

	if collectionName == constants.StudentCollection || collectionName == constants.TeacherCollection {
		err := UpdateSession[T](db, rdb, collectionName, id.Hex(), pipeline)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateSession[T any](
	db *mongo.Database,
	rdb *redis.Client,
	collectionName string,
	userID string,
	pipeline []bson.M,
) error {
	result, err := rdb.Get(context.Background(), fmt.Sprintf(constants.UserKeyStart, userID)).Result()
	if err == redis.Nil {
		return nil
	}

	if result == "" {
		return nil
	}

	repo := mongodb.NewRepository[struct{}, T](db.Collection(collectionName))
	agg, err := repo.AggregateAll(context.Background(), pipeline)
	if err != nil {
		return err
	}

	userData, _ := json.Marshal(agg[0])
	rdb.Set(context.Background(), fmt.Sprintf(constants.UserKeyStart, userID), userData, redis.KeepTTL)

	return nil
}

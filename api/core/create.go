package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ms-admin/api/constants"
	"ms-admin/api/utils"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReferenceCheck struct {
	Collection string
	ID         primitive.ObjectID
	ErrMsg     string
}

func CreateDocument[T ValidatableModel](
	db *mongo.Database,
	data json.RawMessage,
	collectionName string,
	uniqueFields bson.M,
	checkReferences []ReferenceCheck,
) (interface{}, error) {
	var model T
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&model); err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidInput)
	}

	if validator, ok := any(model).(ValidatableModel); ok {
		if err := validator.Validate(); err != nil {
			return nil, err
		}
	}

	if len(uniqueFields) > 0 {
		if err := utils.CheckReplica(db, collectionName, uniqueFields); err != nil {
			return nil, err
		}
	}

	for _, ref := range checkReferences {
		repo := mongodb.NewRepository[any, struct{}](db.Collection(ref.Collection))
		_, err := repo.FindOne(context.Background(), bson.M{"_id": ref.ID})
		if err != nil {
			return nil, fmt.Errorf(ref.ErrMsg)
		}
	}

	repo := mongodb.NewRepository[T, struct{}](db.Collection(collectionName))
	docID, err := repo.InsertOne(context.Background(), &model)
	if err != nil {
		return nil, err
	}

	return docID, nil
}

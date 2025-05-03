package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ms-admin/api/messages"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
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
		return nil, fmt.Errorf("%v: %v", messages.ErrInvalidDataStatus, err)
	}

	fields := map[string]string{
		"status": status.Status,
	}

	for name, value := range fields {
		if value == "" {
			return nil, fmt.Errorf(messages.ErrFieldCannotEmpty, name)
		}
	}

	err := CheckReplica(db, StatusCollection, bson.M{"status": status.Status})
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	statusRepo := mongodb.NewRepository[models.StatusesModel, interface{}](db.Collection(StatusCollection))
	statusID, err := statusRepo.InsertOne(context.Background(), &status)
	if err != nil {
		return nil, err
	}

	return statusID, nil
}

func ReadStatuses(db *mongo.Database) (interface{}, []string, error) {
	statusRepo := mongodb.NewRepository[models.StatusesModel, interface{}](db.Collection(StatusCollection))
	statusFind, err := statusRepo.FindAll(context.Background(), bson.M{})
	if err != nil {
		return nil, nil, err
	}

	var structForHead models.StatusesModel
	header := GetFieldNames(structForHead)
	header = FilterHeaders(header, []string{"ID"})

	return statusFind, header, nil
}

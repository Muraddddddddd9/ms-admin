package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"ms-admin/api/constants"
	"strings"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateObjects(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var object models.ObjectsModel

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&object); err != nil {
		return nil, fmt.Errorf("%v: %v", constants.ErrInvalidDataObject, err)
	}

	object.Object = strings.TrimSpace(strings.ToLower(object.Object))

	fileds := map[string]string{
		"object": object.Object,
	}

	for name, value := range fileds {
		if value == "" {
			return nil, fmt.Errorf(constants.ErrFieldCannotEmpty, name)
		}
	}

	err := CheckReplica(db, constants.ObjectCollection, bson.M{"object": object.Object})
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	objectRepo := mongodb.NewRepository[models.ObjectsModel, struct{}](db.Collection(constants.ObjectCollection))
	objectID, err := objectRepo.InsertOne(context.Background(), &object)
	if err != nil {
		return nil, err
	}

	return objectID, nil
}

func ReadObjects(db *mongo.Database) (interface{}, []string, error) {
	objectRepo := mongodb.NewRepository[models.ObjectsModel, struct{}](db.Collection(constants.ObjectCollection))
	objectFind, err := objectRepo.FindAll(context.Background(), bson.M{})
	if err != nil {
		return nil, nil, err
	}

	var structForHead models.ObjectsModel
	header := GetFieldNames(structForHead)
	header = FilterHeaders(header, []string{"ID"})

	return objectFind, header, nil
}

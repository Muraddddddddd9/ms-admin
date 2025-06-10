package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ms-admin/api/constants"
	"ms-admin/api/core"
	"ms-admin/api/utils"
	"strings"

	"github.com/Muraddddddddd9/ms-database/models"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateObjects(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var object models.ObjectsModel
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&object); err != nil {
		return nil, fmt.Errorf("%v: %v", constants.ErrInvalidDataObject, err)
	}

	return core.CreateDocument[*core.ObjectsModel](
		db,
		data,
		constants.ObjectCollection,
		bson.M{"object": strings.TrimSpace(strings.ToLower(object.Object))},
		nil,
	)
}

func ReadObjects(db *mongo.Database, page, pageSize int) (map[string]interface{}, error) {
	return core.ReadFindDocument[models.ObjectsModel](
		db,
		constants.ObjectCollection,
		[]string{"ID"},
		page, pageSize,
	)
}

func DeleteObjects(db *mongo.Database, collectionName string, ids []primitive.ObjectID) (string, error) {
	checkReferencesOther := []core.ReferenceCheckOther{
		{Collection: constants.ObjectGroupCollection, Field: "object"},
	}

	return core.DeleteDocument[models.ObjectsModel](
		db,
		collectionName,
		ids,
		checkReferencesOther,
	)
}

func UpdateObjects(
	db *mongo.Database,
	rdb *redis.Client,
	collection string,
	id primitive.ObjectID,
	label string,
	newData string,
) error {
	newData = strings.ToLower(newData)
	if label == "object" {
		err := utils.CheckReplica(db, constants.ObjectCollection, bson.M{label: newData})
		if err != nil {
			return err
		}
	}

	return core.UpdateDocument[models.ObjectsModel, struct{}](
		db,
		rdb,
		id,
		constants.ObjectCollection,
		label,
		newData,
		nil,
		nil,
	)
}

package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ms-admin/api/constants"
	"ms-admin/api/core"
	"strings"

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

	return core.CreateDocument[*core.ObjectsModel](
		db,
		data,
		constants.ObjectCollection,
		bson.M{"object": strings.TrimSpace(strings.ToLower(object.Object))},
		nil,
	)
}

func ReadObjects(db *mongo.Database) (map[string]interface{}, error) {
	return core.ReadFindDocument[models.ObjectsModel](
		db,
		constants.ObjectCollection,
		[]string{"ID"},
	)
}

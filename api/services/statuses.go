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

func CreateStatuses(db *mongo.Database, data json.RawMessage) (interface{}, error) {
	var status models.StatusesModel
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&status); err != nil {
		return nil, fmt.Errorf("%v: %v", constants.ErrInvalidDataStatus, err)
	}

	return core.CreateDocument[*core.StatusesModel](
		db,
		data,
		constants.StatusCollection,
		bson.M{"status": strings.TrimSpace(strings.ToLower(status.Status))},
		nil,
	)
}

func ReadStatuses(db *mongo.Database) (map[string]interface{}, error) {
	return core.ReadFindDocument[models.StatusesModel](
		db,
		constants.StatusCollection,
		[]string{"ID"},
	)
}

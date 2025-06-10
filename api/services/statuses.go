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

func ReadStatuses(db *mongo.Database, page, pageSize int) (map[string]interface{}, error) {
	return core.ReadFindDocument[models.StatusesModel](
		db,
		constants.StatusCollection,
		[]string{"ID"},
		page, pageSize,
	)
}

func DeleteStatuses(db *mongo.Database, collectionName string, ids []primitive.ObjectID) (string, error) {
	checkReferencesOther := []core.ReferenceCheckOther{
		{Collection: constants.StudentCollection, Field: "status"},
		{Collection: constants.TeacherCollection, Field: "status"},
	}

	return core.DeleteDocument[models.ObjectsModel](
		db,
		collectionName,
		ids,
		checkReferencesOther,
	)
}

func UpdateStatuses(
	db *mongo.Database,
	rdb *redis.Client,
	collection string,
	id primitive.ObjectID,
	label string,
	newData string,
) error {
	if label == "status" {
		err := utils.CheckReplica(db, constants.StatusCollection, bson.M{label: newData})
		if err != nil {
			return err
		}
	}

	return core.UpdateDocument[models.StatusesModel, struct{}](
		db,
		rdb,
		id,
		constants.StatusCollection,
		label,
		newData,
		nil,
		nil,
	)
}

package core

import (
	"context"
	"fmt"
	"ms-admin/api/constants"
	"ms-admin/api/utils"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReferenceCheckOther struct {
	Collection string
	Field      string
}

func DeleteDocument[T any](
	db *mongo.Database,
	collectionName string,
	ids []primitive.ObjectID,
	checkReferencesOther []ReferenceCheckOther,
) (string, error) {
	var noneDelete = len(ids)
	var flag = true
	for _, id := range ids {
		for _, ref := range checkReferencesOther {
			if err := utils.CheckDataOtherTable(db, ref.Collection, bson.M{ref.Field: id}); err != nil {
				noneDelete--
				flag = false
				break
			}
			flag = true
		}

		if flag {
			repo := mongodb.NewRepository[T, struct{}](db.Collection(collectionName))
			err := repo.DeleteOne(context.Background(), bson.M{"_id": id})
			if err != nil {
				return "", fmt.Errorf(constants.ErrDeleteData)
			}
		}
	}

	if noneDelete != len(ids) {
		return "", fmt.Errorf(constants.SuccDataDelete, noneDelete, len(ids))
	}

	return fmt.Sprintf(constants.SuccDataDelete, noneDelete, len(ids)), nil
}

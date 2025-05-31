package core

import (
	"context"
	"fmt"
	"ms-admin/api/utils"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ReadAggregateDocument[T any](
	db *mongo.Database,
	collectionName string,
	pipeline []bson.M,
	filterHead []string,
	arrSelect []string,
) (map[string]interface{}, error) {
	repo := mongodb.NewRepository[struct{}, T](db.Collection(collectionName))
	agg, err := repo.AggregateAll(context.Background(), pipeline)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	var selRes models.SelectModels
	for _, sel := range arrSelect {
		err := utils.SelectData(db, sel, &selRes)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
	}

	var structForHead T
	header := utils.GetFieldNames(structForHead)
	header = utils.FilterHeaders(header, filterHead)

	arrResult := map[string]interface{}{
		"data":       agg,
		"header":     header,
		"selectData": selRes,
	}

	return arrResult, nil
}

func ReadFindDocument[T any](
	db *mongo.Database,
	collectionName string,
	filterHead []string,
) (map[string]interface{}, error) {
	repo := mongodb.NewRepository[T, struct{}](db.Collection(collectionName))
	findAll, err := repo.FindAll(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	var structForHead T
	header := utils.GetFieldNames(structForHead)
	header = utils.FilterHeaders(header, filterHead)

	arrResult := map[string]interface{}{
		"data":   findAll,
		"header": header,
	}

	return arrResult, nil
}

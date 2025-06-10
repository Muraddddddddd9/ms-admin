package core

import (
	"context"
	"fmt"
	"ms-admin/api/utils"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ReadAggregateDocument[T any](
	db *mongo.Database,
	collectionName string,
	pipeline []bson.M,
	filterHead []string,
	arrSelect []string,
	page, pageSize int,
) (map[string]interface{}, error) {

	if page > 0 && pageSize > 0 {
		skip := (page - 1) * pageSize
		paginationStage := []bson.M{
			{"$skip": skip},
			{"$limit": pageSize},
		}
		pipeline = append(pipeline, paginationStage...)
	}
	
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

	countAll, err := repo.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	arrResult := map[string]interface{}{
		"data":       agg,
		"header":     header,
		"selectData": selRes,
		"count":      countAll,
	}

	return arrResult, nil
}

func ReadFindDocument[T any](
	db *mongo.Database,
	collectionName string,
	filterHead []string,
	page, pageSize int,
) (map[string]interface{}, error) {
	var optionsFind *options.FindOptions
	if page > 0 && pageSize > 0 {
		skip := int64((page - 1) * pageSize)
		limit := int64(pageSize)

		optionsFind = &options.FindOptions{
			Skip:  &skip,
			Limit: &limit,
		}
	} else {
		optionsFind = &options.FindOptions{}
	}

	repo := mongodb.NewRepository[T, struct{}](db.Collection(collectionName))
	findAll, err := repo.FindAll(context.Background(), bson.M{}, optionsFind)
	if err != nil {
		return nil, err
	}

	countAll, err := repo.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	var structForHead T
	header := utils.GetFieldNames(structForHead)
	header = utils.FilterHeaders(header, filterHead)

	arrResult := map[string]interface{}{
		"data":   findAll,
		"header": header,
		"count":  countAll,
	}

	return arrResult, nil
}

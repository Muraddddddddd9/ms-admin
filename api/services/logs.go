package services

import (
	"context"
	"ms-admin/api/constants"
	"ms-admin/api/utils"
	"time"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Logging(db *mongo.Database, api, method, status string, data any, errData any) {
	document := models.Log{
		API:    api,
		Method: method,
		Status: status,
		Data:   data,
		Date:   time.Now().Local().Format("2006-01-02 15:04:05 MST"),
		Error:  errData,
	}
	logRepo := mongodb.NewRepository[models.Log, struct{}](db.Collection(constants.LogsCollection))
	_, err := logRepo.InsertOne(context.Background(), &document)
	if err != nil {
		log.Errorf(constants.ErrDataLogging)
	}
}

func ReadLogs(db *mongo.Database) (map[string]interface{}, error) {
	logRepo := mongodb.NewRepository[models.Log, struct{}](db.Collection(constants.LogsCollection))
	sortOpts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})
	logFind, err := logRepo.FindAll(context.Background(), bson.M{}, sortOpts)
	if err != nil {
		return nil, err
	}

	var structForHead models.Log
	header := utils.GetFieldNames(structForHead)

	arrResult := map[string]interface{}{
		"data":   logFind,
		"header": header,
	}

	return arrResult, nil
}

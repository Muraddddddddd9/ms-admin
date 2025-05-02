package services

import (
	"context"
	"time"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	LogsCollection = "log"
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
	logRepo := mongodb.NewRepository[models.Log, interface{}](db.Collection(LogsCollection))
	_, err := logRepo.InsertOne(context.Background(), &document)
	if err != nil {
		log.Errorf("Ошибка в логгирование данных")
	}
}

func ReadLogs(db *mongo.Database) (interface{}, []string, error) {
	logRepo := mongodb.NewRepository[models.Log, interface{}](db.Collection(LogsCollection))

	sortOpts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})
	logFind, err := logRepo.FindAll(context.Background(), bson.M{}, sortOpts)
	if err != nil {
		return nil, nil, err
	}

	var structForHead models.Log
	header := GetFieldNames(structForHead)
	header = FilterHeaders(header, []string{"ID"})

	return logFind, header, nil
}

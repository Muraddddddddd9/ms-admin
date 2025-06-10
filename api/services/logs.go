package services

import (
	"context"
	"encoding/json"
	"ms-admin/api/constants"
	"ms-admin/api/utils"
	"strconv"
	"strings"
	"time"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"github.com/gofiber/fiber/v2"
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

func ReadLogs(c *fiber.Ctx, db *mongo.Database) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	skip := (page - 1) * limit

	filter := bson.M{}

	if apiParam := c.Query("api"); apiParam != "" {
		apis := strings.Split(apiParam, ",")
		filter["api"] = bson.M{"$in": apis}
	}

	if methodParam := c.Query("method"); methodParam != "" {
		methods := strings.Split(methodParam, ",")
		filter["method"] = bson.M{"$in": methods}
	}

	if statusParam := c.Query("status"); statusParam != "" {
		var statusFilters []struct {
			Operator string `json:"operator"`
			Value    struct {
				Start int `json:"start"`
				End   int `json:"end"`
			} `json:"value"`
		}

		if err := json.Unmarshal([]byte(statusParam), &statusFilters); err == nil {
			var statusConditions []bson.M
			for _, cond := range statusFilters {
				if cond.Operator == "range" {
					startStr := strconv.Itoa(cond.Value.Start)
					endStr := strconv.Itoa(cond.Value.End)

					statusConditions = append(statusConditions, bson.M{
						"status": bson.M{
							"$gte": startStr,
							"$lt":  endStr,
						},
					})
				}
			}

			if len(statusConditions) > 0 {
				filter["$or"] = statusConditions
			}
		}
	}

	logRepo := mongodb.NewRepository[models.Log, struct{}](db.Collection(constants.LogsCollection))

	total, err := logRepo.CountDocuments(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	sortOpts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	logFind, err := logRepo.FindAll(context.Background(), filter, sortOpts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var structForHead models.Log
	header := utils.GetFieldNames(structForHead)
	header = utils.FilterHeaders(header, []string{"ID"})

	arrResult := map[string]interface{}{
		"data":    logFind,
		"header":  header,
		"hasMore": int64((page * limit)) < total,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": arrResult,
	})
}

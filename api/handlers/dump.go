package handlers

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"ms-admin/api/constants"
	"ms-admin/api/services"
	loconfig "ms-admin/config"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionStruct struct {
	Collections []string `json:"collections"`
}

func GetCollections(c *fiber.Ctx, db *mongo.Database) error {
	colletions, err := db.ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf(constants.ErrCollectionsNotFound, err),
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"data": colletions,
	})
}

func Dump(c *fiber.Ctx, db *mongo.Database) error {
	var collectionsData CollectionStruct

	if err := c.BodyParser(&collectionsData); err != nil {
		services.Logging(db, "/api/admin/dump", c.Method(), strconv.Itoa(fiber.StatusBadRequest), collectionsData, err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrInvalidInput,
		})
	}

	if len(collectionsData.Collections) == 0 {
		services.Logging(db, "/api/admin/dump", c.Method(), strconv.Itoa(fiber.StatusConflict), collectionsData, constants.ErrInputCollection)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": constants.ErrInputCollection,
		})
	}

	now := time.Now()
	date := now.Local().Format("2006-01-02")
	hour := fmt.Sprintf("%02d", now.Hour())
	min := fmt.Sprintf("%02d", now.Minute())
	sec := fmt.Sprintf("%02d", now.Second())
	zipName := fmt.Sprintf("mongodump_%v_%v-%v-%v.zip", date, hour, min, sec)

	pr, pw := io.Pipe()
	zipWriter := zip.NewWriter(pw)
	cfg, err := loconfig.LoadLocalConfig()
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": constants.ErrLoadEnv,
		})
	}

	go func() {
		defer pw.Close()
		defer zipWriter.Close()

		for _, collection := range collectionsData.Collections {
			if collection == "system.views" || collection == "system.profile" {
				continue
			}

			cursor, err := db.Collection(collection).Find(context.Background(), bson.M{"status": bson.M{"$ne": constants.AdminStatus}, "email": bson.M{"$ne": cfg.ADMIN_EMAIL}})
			if err != nil {
				continue
			}
			defer cursor.Close(context.Background())

			fileInZip, err := zipWriter.Create(collection + ".json")
			if err != nil {
				continue
			}

			for cursor.Next(context.Background()) {
				var doc bson.M
				if err := cursor.Decode(&doc); err != nil {
					continue
				}

				jsonData, err := json.Marshal(doc)
				if err != nil {
					continue
				}

				if _, err := fileInZip.Write(jsonData); err != nil {
					continue
				}

				fileInZip.Write([]byte("\n"))
			}
		}
	}()

	c.Set("Content-Type", "application/zip")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, zipName))

	return c.SendStream(pr)
}

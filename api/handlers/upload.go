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
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/Muraddddddddd9/ms-database/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var handlerCollection = map[string]interface{}{
	constants.StudentCollection:     models.StudentsModel{},
	constants.TeacherCollection:     models.TeachersModel{},
	constants.GroupCollection:       models.GroupsModel{},
	constants.ObjectCollection:      models.ObjectsModel{},
	constants.ObjectGroupCollection: models.ObjectsGroupsModel{},
	constants.StatusCollection:      models.StatusesModel{},
}

func Upload(c *fiber.Ctx, db *mongo.Database) error {
	fileZip, err := c.FormFile("upload")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrLoadFile,
		})
	}
	type FileData struct {
		File string `bson:"file"`
	}

	fileData := FileData{File: fileZip.Filename}

	file, err := fileZip.Open()
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	tempFile, err := os.CreateTemp("", fmt.Sprintf("%v-*", fileZip.Filename)+filepath.Ext(fileZip.Filename))
	if err != nil {
		file.Close()
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	tempFilePath := tempFile.Name()
	defer os.Remove(tempFilePath)

	_, err = io.Copy(tempFile, file)
	file.Close()
	tempFile.Close()
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	files, err := zip.OpenReader(tempFile.Name())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constants.ErrServerError,
		})
	}

	for _, file := range files.File {
		fileName := strings.Split(file.Name, ".")[0]
		model, ok := handlerCollection[fileName]
		if !ok {
			continue
		}

		collection, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": constants.ErrServerError,
			})
		}
		defer collection.Close()

		data, err := io.ReadAll(collection)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": constants.ErrServerError,
			})
		}

		var insertData []interface{}
		stringArr := strings.Split(string(data), "\n")

		for _, line := range stringArr {
			if line == "" {
				continue
			}

			item := reflect.New(reflect.TypeOf(model)).Interface()

			if err := json.Unmarshal([]byte(line), &item); err != nil {
				continue
			}

			insertData = append(insertData, item)
		}

		if len(insertData) > 0 {
			cfg, err := loconfig.LoadLocalConfig()
			if err != nil {
				return err
			}

			_, err = db.Collection(fileName).DeleteMany(context.Background(), bson.M{"status": bson.M{"$ne": constants.AdminStatus}, "email": bson.M{"$ne": cfg.ADMIN_EMAIL}})
			if err != nil {
				services.Logging(db, "/api/admin/upload", c.Method(), strconv.Itoa(fiber.StatusBadRequest), fileData, err.Error())
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": constants.ErrDeleteData,
				})
			}
			_, err = db.Collection(fileName).InsertMany(context.Background(), insertData)
			if err != nil {
				services.Logging(db, "/api/admin/upload", c.Method(), strconv.Itoa(fiber.StatusBadRequest), fileData, err.Error())
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": constants.ErrCreateData,
				})
			}

		}
	}

	services.Logging(db, "/api/admin/upload", c.Method(), strconv.Itoa(fiber.StatusAccepted), fileData, nil)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": constants.SuccUploadFile,
	})
}

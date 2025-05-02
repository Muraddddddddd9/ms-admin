package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	loconfig "ms-admin/config"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func CreateStartAdmin(db *mongo.Database, cfg *loconfig.LocalConfig) error {
	if cfg.ADMIN_EMAIL == "" || cfg.ADMIN_PASSWORD == "" {
		return errors.New("admin email/password not set in config")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(cfg.ADMIN_PASSWORD),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	teacherRepo := mongodb.NewRepository[models.TeachersModel, models.TeachersWithStatusModel](db.Collection("teachers"))
	statuseRepo := mongodb.NewRepository[models.StatusesModel, interface{}](db.Collection("statuses"))

	var newAdminId interface{}
	if _, err = statuseRepo.FindOne(context.Background(), bson.M{"status": "админ"}); err != nil {
		if err == mongo.ErrNoDocuments {
			newAdminId, err = statuseRepo.InsertOne(context.Background(), &models.StatusesModel{Status: "админ"})
			if err != nil {
				return fmt.Errorf("failed to create status: %v", err)
			}
		} else {
			return fmt.Errorf("failed to find status: %v", err)
		}
	}

	if _, err := teacherRepo.FindOne(context.Background(), bson.M{"email": cfg.ADMIN_EMAIL}); err != nil {
		if err == mongo.ErrNoDocuments {
			document := models.TeachersModel{
				Name:       "Admin",
				Surname:    "Admin",
				Patronymic: "Admin",
				Email:      cfg.ADMIN_EMAIL,
				Password:   string(hashedPassword),
				Status:     newAdminId.(primitive.ObjectID),
			}

			_, err = teacherRepo.InsertOne(context.Background(), &document)
			if err != nil {
				return fmt.Errorf("failed to create admin: %v", err)
			}
			log.Println("Администратор создан")
		} else {
			return fmt.Errorf("failed to check admin existence: %v", err)
		}
	} else {
		log.Println("Администратор уже существует")
	}

	return nil
}

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"ms-admin/api/constants"
	loconfig "ms-admin/config"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func CreateAdmin(db *mongo.Database, cfg *loconfig.LocalConfig, statusID primitive.ObjectID) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(cfg.ADMIN_PASSWORD),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf(constants.ErrHashPassword, err)
	}

	teacherRepo := mongodb.NewRepository[models.TeachersModel, models.TeachersWithStatusModel](db.Collection(constants.TeacherCollection))
	if _, err := teacherRepo.FindOne(context.Background(), bson.M{"email": cfg.ADMIN_EMAIL}); err != nil {
		if err == mongo.ErrNoDocuments {
			document := models.TeachersModel{
				Name:       "Admin",
				Surname:    "Admin",
				Patronymic: "Admin",
				Email:      cfg.ADMIN_EMAIL,
				Password:   string(hashedPassword),
				IPs:        []string{},
				Status:     statusID,
			}

			_, err = teacherRepo.InsertOne(context.Background(), &document)
			if err != nil {
				return fmt.Errorf(constants.ErrCreateAdmin, err)
			}
			log.Println(constants.SuccCreateAdmin)
		} else {
			return fmt.Errorf(constants.ErrCheckAdmin, err)
		}
	} else {
		log.Print(constants.SuccAdminAlreadyExists)
	}

	return nil
}

func CreateStartAdmin(db *mongo.Database, cfg *loconfig.LocalConfig) error {
	if cfg.ADMIN_EMAIL == "" || cfg.ADMIN_PASSWORD == "" {
		return errors.New(constants.ErrAdminConfig)
	}

	var newAdminId interface{}
	statuseRepo := mongodb.NewRepository[models.StatusesModel, interface{}](db.Collection(constants.StatusCollection))
	if adminID, err := statuseRepo.FindOne(context.Background(), bson.M{"status": constants.AdminCreate}); err != nil {
		if err == mongo.ErrNoDocuments {
			newAdminId, err = statuseRepo.InsertOne(context.Background(), &models.StatusesModel{Status: constants.AdminCreate})
			if err != nil {
				return fmt.Errorf(constants.ErrCreateAdminStatus, err)
			}
			err = CreateAdmin(db, cfg, newAdminId.(primitive.ObjectID))
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf(constants.ErrAdminNotFound, err)
		}
	} else {
		err = CreateAdmin(db, cfg, adminID.ID)
		if err != nil {
			return fmt.Errorf(constants.ErrCreateAdminStatus, err)
		}
	}

	return nil
}

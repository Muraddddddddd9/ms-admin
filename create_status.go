package main

import (
	"encoding/json"
	"log"
	"ms-admin/api/constants"
	"ms-admin/api/services"

	"go.mongodb.org/mongo-driver/mongo"
)

func CreateStatusAll(db *mongo.Database) {
	arrStatus := []string{constants.RestrictedAdminStatusCreate, constants.TeacherStatusCreate, constants.StudentStatusCreate}
	for _, v := range arrStatus {
		input := map[string]string{"status": v}
		byteV, err := json.Marshal(input)
		if err != nil {
			log.Printf(constants.ErrCreateStatus, v, err.Error())
			continue
		}

		id, err := services.CreateStatuses(db, json.RawMessage(byteV))
		if err != nil {
			log.Printf(constants.ErrCreateStatus, v, err.Error())
			continue
		}

		log.Printf("Стауст %v был создан с ID = %v", v, id)
	}
}

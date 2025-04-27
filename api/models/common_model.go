package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SelectModels struct {
	Objects  []SelectObjectModels  `json:"select_object"`
	Groups   []SelectGroupModels   `json:"select_groups"`
	Teachers []SelectTeacherModels `json:"select_teachers"`
	Statuses []SelectStatusModels  `json:"select_statuses"`
}

type SelectObjectModels struct {
	ID     primitive.ObjectID `bson:"_id" json:"value"`
	Object string             `bson:"object" json:"label"`
}

type SelectGroupModels struct {
	ID    primitive.ObjectID `bson:"_id" json:"value"`
	Group string             `bson:"group" json:"label"`
}

type SelectStatusModels struct {
	ID     primitive.ObjectID `bson:"_id" json:"value"`
	Status string             `bson:"status" json:"label"`
}

type SelectTeacherModels struct {
	ID         primitive.ObjectID `bson:"_id" json:"value"`
	Name       string             `bson:"name" json:"name"`
	Surname    string             `bson:"surname" json:"surname"`
	Patronymic string             `bson:"patronymic" json:"patronymic"`
}

type Log struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	API    string             `bson:"api" json:"api"`
	Method string             `bson:"method" json:"method"`
	Status string             `bson:"status" json:"status"`
	Data   any                `bson:"data" json:"data"`
	Date   string             `bson:"date" json:"date"`
	Error  any                `bson:"error" json:"error"`
}

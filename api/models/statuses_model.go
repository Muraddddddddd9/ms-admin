package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type StatusesModel struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Status string             `bson:"status" json:"status"`
}

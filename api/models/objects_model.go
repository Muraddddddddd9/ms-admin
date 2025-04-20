package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ObjectsModel struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Object string             `bson:"object" json:"object"`
}

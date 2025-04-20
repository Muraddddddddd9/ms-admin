package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type GroupsModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Group     string             `bson:"group" json:"group"`
	TeacherId primitive.ObjectID `bson:"teacher_id" json:"teacher_id"`
}

type GroupWithTeacher struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Group     string             `bson:"group" json:"group"`
	TeacherId string             `bson:"teacher_id" json:"teacher_id"`
}

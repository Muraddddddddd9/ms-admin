package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type GroupsModel struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Group   string             `bson:"group" json:"group"`
	Teacher primitive.ObjectID `bson:"teacher" json:"teacher"`
}

type GroupsWithTeacherModel struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Group   string             `bson:"group" json:"group"`
	Teacher string             `bson:"teacher" json:"teacher"`
}

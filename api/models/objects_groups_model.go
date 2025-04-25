package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ObjectsGroupsModel struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Object  primitive.ObjectID `bson:"object" json:"object"`
	Group   primitive.ObjectID `bson:"group" json:"group"`
	Teacher primitive.ObjectID `bson:"teacher" json:"teacher"`
}

type ObjectsGroupsWithGroupAndTeacherModel struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Object  string             `bson:"object" json:"object"`
	Group   string             `bson:"group" json:"group"`
	Teacher string             `bson:"teacher" json:"teacher"`
}

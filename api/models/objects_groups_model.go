package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ObjectsGroupsModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ObjectId  primitive.ObjectID `bson:"object_id" json:"object_id"`
	GroupId   primitive.ObjectID `bson:"group_id" json:"group_id"`
	TeacherId primitive.ObjectID `bson:"teacher_id" json:"teacher_id"`
}

type ObjectsGroupsWithGroupAndTeacherModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ObjectId  string             `bson:"object_id" json:"object_id"`
	GroupId   string             `bson:"group_id" json:"group_id"`
	TeacherId string             `bson:"teacher_id" json:"teacher_id"`
}

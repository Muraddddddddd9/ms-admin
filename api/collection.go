package api

import "go.mongodb.org/mongo-driver/bson/primitive"

type Students struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	Surname    string             `bson:"surname"`
	Patronymic string             `bson:"patronymic"`
	Group      primitive.ObjectID `bson:"group"`
	Email      string             `bson:"email"`
	Password   string             `bson:"password"`
	Telegram   string             `bson:"telegram"`
	Diplomas   []string           `bson:"diplomas"`
	Ips        []string           `bson:"ips"`
	Status     string             `bson:"status"`
}

type Teachers struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	Surname    string             `bson:"surname"`
	Patronymic string             `bson:"patronymic"`
	Email      string             `bson:"email"`
	Password   string             `bson:"password"`
	Telegram   string             `bson:"telegram"`
	Ips        []string           `bson:"ips"`
	Status     string             `bson:"status"`
}

type Groups struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Group      string             `bson:"group"`
	TeachersId primitive.ObjectID `bson:"teachers_id"`
}

type Objects struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Object string             `bson:"object"`
}

type ObjectsForGroups struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ObjectId   primitive.ObjectID `json:"object_id"`
	GroupId    primitive.ObjectID `json:"group_id"`
	TeachersId primitive.ObjectID `json:"teachers_id"`
}

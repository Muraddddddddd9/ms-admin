package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type StudentsModel struct {
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

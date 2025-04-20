package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type StudentsModel struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name       string             `bson:"name" json:"name"`
	Surname    string             `bson:"surname" json:"surname"`
	Patronymic string             `bson:"patronymic" json:"patronymic"`
	Group      primitive.ObjectID `bson:"group" json:"group"`
	Email      string             `bson:"email" json:"email"`
	Password   string             `bson:"password" json:"password"`
	Telegram   string             `bson:"telegram" json:"telegram"`
	Diplomas   []string           `bson:"diplomas" json:"diplomas"`
	IPs        []string           `bson:"ips" json:"ips"`
	Status     string             `bson:"status" json:"status"`
}

type StudentWithGroup struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name       string             `bson:"name" json:"name"`
	Surname    string             `bson:"surname" json:"surname"`
	Patronymic string             `bson:"patronymic" json:"patronymic"`
	Group      string             `bson:"group" json:"group"`
	Email      string             `bson:"email" json:"email"`
	Password   string             `bson:"password" json:"password"`
	Telegram   string             `bson:"telegram" json:"telegram"`
	Diplomas   []string           `bson:"diplomas" json:"diplomas"`
	IPs        []string           `bson:"ips" json:"ips"`
	Status     string             `bson:"status" json:"status"`
}

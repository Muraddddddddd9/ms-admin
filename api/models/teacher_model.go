package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TeachersModel struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name       string             `bson:"name" json:"name"`
	Surname    string             `bson:"surname" json:"surname"`
	Patronymic string             `bson:"patronymic" json:"patronymic"`
	Email      string             `bson:"email" json:"email"`
	Password   string             `bson:"password" json:"password"`
	Telegram   string             `bson:"telegram" json:"telegram"`
	IPs        []string           `bson:"ips" json:"ips"`
	Status     primitive.ObjectID `bson:"status" json:"status"`
}

type TeachersWithStatusModel struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name       string             `bson:"name" json:"name"`
	Surname    string             `bson:"surname" json:"surname"`
	Patronymic string             `bson:"patronymic" json:"patronymic"`
	Email      string             `bson:"email" json:"email"`
	Password   string             `bson:"password" json:"password"`
	Telegram   string             `bson:"telegram" json:"telegram"`
	IPs        []string           `bson:"ips" json:"ips"`
	Status     string             `bson:"status" json:"status"`
}

package data

import (
	"context"
	"fmt"
	"log"
	"ms-admin/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo() (*mongo.Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// mongodb://college:BIM_LOCAL1@localhost:27017/admin?authSource=admin
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/admin?authSource=%s", cfg.DB_USERNAME, cfg.DB_PASSWORD, cfg.DB_HOST, cfg.DB_PORT, cfg.DB_AUTH_SOURCE)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client, nil
}

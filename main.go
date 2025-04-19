package main

import (
	"context"
	"log"
	"ms-admin/api"
	"ms-admin/data"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	client, err := data.ConnectMongo()
	if err == nil {
		log.Println("Connect to MONGODB is succ")
	}
	defer client.Disconnect(context.TODO())

	store, err := data.ConnectRedis()
	if err == nil {
		log.Println("Connect to REDIS is succ")
	}
	defer store.Close()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE",
	}))

	app.Post("/api/admin/create_data", func(c *fiber.Ctx) error {
		return api.CreateData(c, client)
	})

	// app.Patch("/api/admin/update_data", func(c *fiber.Ctx) error {
	// 	return
	// })
}

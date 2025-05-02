package main

import (
	"context"
	"fmt"
	"log"
	"ms-admin/api/handlers"
	"ms-admin/config"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/data/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	db, err := mongodb.Connect()
	if err == nil {
		log.Println("Connect to MONGODB is succ")
	}
	defer db.Client().Disconnect(context.Background())

	storage, err := redis.Connect()
	if err == nil {
		log.Println("Connect to REDIS is succ")
	}
	defer storage.Close()

	cfg, err := loconfig.LoadLocalConfig()
	if err != nil {
		log.Print(err)
	}

	err = CreateStartAdmin(db, cfg)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	app.Post("/api/admin/create_data", AdminOnly(storage), func(c *fiber.Ctx) error {
		return handlers.CreateData(c, db)
	})

	app.Get("/api/admin/get_data/:collection", AdminOnly(storage), func(c *fiber.Ctx) error {
		return handlers.GetData(c, db)
	})

	app.Patch("/api/admin/update_data", AdminOnly(storage), func(c *fiber.Ctx) error {
		return handlers.UpdateData(c, db)
	})

	app.Delete("/api/admin/delete_data", AdminOnly(storage), func(c *fiber.Ctx) error {
		return handlers.DeleteData(c, db)
	})

	app.Listen(fmt.Sprintf("localhost%v", cfg.PROJECT_PORT))
}

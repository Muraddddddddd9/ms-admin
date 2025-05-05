package main

import (
	"context"
	"fmt"
	"log"
	"ms-admin/api/constants"
	"ms-admin/api/handlers"
	loconfig "ms-admin/config"

	"github.com/Muraddddddddd9/ms-database/data/mongodb"
	"github.com/Muraddddddddd9/ms-database/data/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	db, err := mongodb.Connect()
	if err == nil {
		log.Println(constants.SuccConnectMongo)
	}
	defer db.Client().Disconnect(context.Background())

	rdb, err := redis.Connect()
	if err == nil {
		log.Println(constants.SuccConnectRedis)
	}
	defer rdb.Close()

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
		AllowOrigins:     cfg.ORIGIN_URL,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	app.Post("/api/admin/create_data", AdminOnly(rdb), func(c *fiber.Ctx) error {
		return handlers.CreateData(c, db)
	})

	app.Get("/api/admin/get_data/:collection", AdminOnly(rdb), func(c *fiber.Ctx) error {
		return handlers.GetData(c, db)
	})

	app.Patch("/api/admin/update_data", AdminOnly(rdb), func(c *fiber.Ctx) error {
		return handlers.UpdateData(c, db, rdb)
	})

	app.Delete("/api/admin/delete_data", AdminOnly(rdb), func(c *fiber.Ctx) error {
		return handlers.DeleteData(c, db)
	})

	// app.Listen(fmt.Sprintf("localhost%v", cfg.PROJECT_PORT))
	app.Listen(fmt.Sprintf("%v", cfg.PROJECT_PORT))
}

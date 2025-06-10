package main

import (
	"context"
	"fmt"
	"log"
	"ms-admin/api/constants"
	"ms-admin/api/handlers"
	"ms-admin/api/services"
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

	CreateStatusAll(db)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.ORIGIN_URL,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	app.Post("/api/admin/create_data", Access(rdb, []string{constants.AdminStatus}), func(c *fiber.Ctx) error {
		return handlers.CreateData(c, db)
	})

	app.Get("/api/admin/get_data/:collection", Access(rdb, []string{constants.AdminStatus}), func(c *fiber.Ctx) error {
		return handlers.GetData(c, db)
	})

	app.Get("/api/admin/get_logs", Access(rdb, []string{constants.AdminStatus}), func(c *fiber.Ctx) error {
		return services.ReadLogs(c, db)
	})

	app.Patch("/api/admin/update_data", Access(rdb, []string{constants.AdminStatus}), func(c *fiber.Ctx) error {
		return handlers.UpdateData(c, db, rdb)
	})

	app.Delete("/api/admin/delete_data", Access(rdb, []string{constants.AdminStatus}), func(c *fiber.Ctx) error {
		return handlers.DeleteData(c, db)
	})

	app.Get("/api/admin/get_all_object/:group", Access(rdb, []string{constants.AdminStatus, constants.RestrictedAdminStatus}), func(c *fiber.Ctx) error {
		return services.GetAllObject(c, db)
	})

	app.Get("/api/admin/get_collections", Access(rdb, []string{constants.AdminStatus}), func(c *fiber.Ctx) error {
		return handlers.GetCollections(c, db)
	})

	app.Post("/api/admin/dump", Access(rdb, []string{constants.AdminStatus}), func(c *fiber.Ctx) error {
		return handlers.Dump(c, db)
	})

	app.Post("/api/admin/drop", Access(rdb, []string{constants.AdminStatus}), func(c *fiber.Ctx) error {
		return handlers.Drop(c, db, cfg)
	})

	app.Patch("/api/admin/up_group", Access(rdb, []string{constants.AdminStatus}), func(c *fiber.Ctx) error {
		return handlers.UpGroup(c, db)
	})

	app.Post("/api/admin/upload", Access(rdb, []string{constants.AdminStatus}), func(c *fiber.Ctx) error {
		return handlers.Upload(c, db)
	})

	go runBackups(db)

	// app.Listen(fmt.Sprintf("localhost%v", cfg.PROJECT_PORT))
	app.Listen(fmt.Sprintf("%v", cfg.PROJECT_PORT))
}

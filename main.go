package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"ms-admin/api/handlers"
	"ms-admin/api/models"
	"ms-admin/config"
	"ms-admin/data"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/storage/redis/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateStartAdmin(db *mongo.Database, cfg *config.Config) error {
	if cfg.ADMIN_EMAIL == "" || cfg.ADMIN_PASSWORD == "" {
		return errors.New("admin email/password not set in config")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(cfg.ADMIN_PASSWORD),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	ctx := context.Background()

	filterStatus := bson.M{"status": "админ"}
	var statusDoc bson.M
	if err := db.Collection("statuses").FindOne(ctx, filterStatus).Decode(&statusDoc); err != nil {
		if err == mongo.ErrNoDocuments {
			res, err := db.Collection("statuses").InsertOne(ctx, bson.M{"status": "админ"})
			if err != nil {
				return fmt.Errorf("failed to create status: %v", err)
			}
			statusDoc = bson.M{"_id": res.InsertedID, "status": "админ"}
		} else {
			return fmt.Errorf("failed to find status: %v", err)
		}
	}

	filterAdmin := bson.M{"email": cfg.ADMIN_EMAIL}
	var existingAdmin bson.M
	if err := db.Collection("teachers").FindOne(ctx, filterAdmin).Decode(&existingAdmin); err != nil {
		if err == mongo.ErrNoDocuments {
			document := models.TeachersModel{
				Name:       "Admin",
				Surname:    "Admin",
				Patronymic: "Admin",
				Email:      cfg.ADMIN_EMAIL,
				Password:   string(hashedPassword),
				Status:     statusDoc["_id"].(primitive.ObjectID),
			}

			_, err = db.Collection("teachers").InsertOne(ctx, document)
			if err != nil {
				return fmt.Errorf("failed to create admin: %v", err)
			}
			log.Println("Администратор создан")
		} else {
			return fmt.Errorf("failed to check admin existence: %v", err)
		}
	} else {
		log.Println("Администратор уже существует")
	}

	return nil
}

func AdminOnly(storage *redis.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session_id := c.Cookies("session_id")

		body, err := storage.Get(fmt.Sprintf("users:%v", session_id))
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Сессия не найдена"})
		}

		var user struct {
			Status string `json:"status"`
		}
		err = json.Unmarshal(body, &user)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Ошибка в получение данных"})
		}

		if user.Status != "админ" {
			return c.Status(404).JSON(fiber.Map{})
		}

		return c.Next()
	}
}

func main() {
	client, err := data.ConnectMongo()
	if err == nil {
		log.Println("Connect to MONGODB is succ")
	}
	defer client.Disconnect(context.TODO())
	db := client.Database("diary")

	storage, err := data.ConnectRedis()
	if err == nil {
		log.Println("Connect to REDIS is succ")
	}
	defer storage.Close()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Print(err)
	}

	CreateStartAdmin(db, cfg)

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

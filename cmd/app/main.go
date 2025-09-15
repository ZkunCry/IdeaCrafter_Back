package main

import (
	"fmt"
	"startup_back/internal/config"
	"startup_back/internal/domain"
	"startup_back/internal/handler"
	"startup_back/internal/repository"
	"startup_back/internal/routes"
	"startup_back/internal/service"

	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}
	logrus.Infof("Loaded config: %+v", cfg)
	db, err := gorm.Open(postgres.Open(cfg.DBConnectionString()),&gorm.Config{})
	if err !=nil{
		logrus.Fatalf("Failed to connect to database: %v", err)

	}
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Startup{},
		&domain.Category{},
		&domain.Membership{},
		&domain.Favorite{},
		&domain.StartupFile{},
	)
	if err != nil{
		logrus.Fatalf("Failed to migrate database: %v", err)
	}
	repos:= repository.NewRespositories(db)
	services:= service.NewServices(repos,&cfg)
	fmt.Printf("Services %v\n",services)
	handlers:=handler.NewHandlers(services)

	app:= fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
				logrus.WithFields(logrus.Fields{
				"error": err.Error(),
				"path":  c.Path(),
			}).Error("Request failed")
			return c.Status(code).JSON(fiber.Map{"error":err.Error()})
		},
	})
	routes.SetupRoutes(app,handlers)
	address := fmt.Sprintf("%s:%s", cfg.Server.Host,cfg.Server.Port)
	logrus.Infof("Starting server on %s", address)
	if err:= app.Listen(address); err!=nil{
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
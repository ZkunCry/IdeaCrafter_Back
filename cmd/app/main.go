package main

import (
	"fmt"
	"startup_back/internal/config"
	"startup_back/internal/domain"
	"startup_back/internal/handler"
	"startup_back/internal/repository"
	"startup_back/internal/routes"
	"startup_back/internal/service"
	"time"

	"github.com/sirupsen/logrus"

	_ "startup_back/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Startup API
// @version 1.0
// @description API для управления стартапами, вакансиями и пользователями
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3001
// @BasePath /api

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
		&domain.Favorite{},
		&domain.StartupFile{},
		&domain.Role{},
		&domain.Vacancy{},
		&domain.Application{},
		&domain.Stage{},
	)

	if err != nil{
		logrus.Fatalf("Failed to migrate database: %v", err)
	}
	repos:= repository.NewRespositories(db)
	services:= service.NewServices(repos,&cfg)

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
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next() 
		duration := time.Since(start)

		entry := logrus.WithFields(logrus.Fields{
			"method":   c.Method(),
			"path":     c.Path(),
			"status":   c.Response().StatusCode(),
			"latency":  duration,
			"clientIP": c.IP(),
		})

		if err != nil {
			entry.WithField("error", err.Error()).Error("Request failed")
		} else {
			entry.Info("Request completed")
		}

		return err
	})
	app.Get("/swagger/*", swagger.HandlerDefault)
	routes.SetupRoutes(app,handlers,services)
	address := fmt.Sprintf("%s:%s", cfg.Server.Host,cfg.Server.Port)
	logrus.Infof("Starting server on %s", address)
	if err:= app.Listen(address); err!=nil{
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
package main

import (
	"fmt"
	"os"
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

var(
	JsonLogger *logrus.Logger
	TermLogger *logrus.Logger

)


func main() {
	cfg, err := config.LoadConfig()
	JsonLogger = logrus.New()
	JsonLogger.SetFormatter(&logrus.JSONFormatter{

	})
	file,err:= os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err!=nil{
		logrus.Fatalf("Failed to open log file: %v", err)
	}
	JsonLogger.SetOutput(file)
	TermLogger = logrus.New()
	TermLogger.SetFormatter(&logrus.TextFormatter{
        TimestampFormat: "2006-01-02 15:04:05",
        FullTimestamp:   true,
        ForceColors:     true,
        DisableLevelTruncation: true,
        PadLevelText:    true,
    })
  TermLogger.SetOutput(os.Stdout)
	if err != nil {
		TermLogger.Fatalf("error loading config: %v", err)
	}
	TermLogger.Infof("Loaded config: %+v", cfg)
	db, err := gorm.Open(postgres.Open(cfg.DBConnectionString()),&gorm.Config{})
	if err !=nil{
		TermLogger.Fatalf("Failed to connect to database: %v", err)

	}
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Startup{},
		&domain.Category{},
		&domain.Membership{},
		&domain.Favorite{},
		&domain.StartupFile{},
		&domain.Role{},
		&domain.Vacancy{},
		&domain.Application{},
	)
	if err != nil{
		TermLogger.Fatalf("Failed to migrate database: %v", err)
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
				JsonLogger.WithFields(logrus.Fields{
				"error": err.Error(),
				"path":  c.Path(),
			}).Error("Request failed")
			return c.Status(code).JSON(fiber.Map{"error":err.Error()})
		},
	})
	routes.SetupRoutes(app,handlers,services)
	address := fmt.Sprintf("%s:%s", cfg.Server.Host,cfg.Server.Port)
	TermLogger.Infof("Starting server on %s", address)
	if err:= app.Listen(address); err!=nil{
		TermLogger.Fatalf("Failed to start server: %v", err)
	}
}
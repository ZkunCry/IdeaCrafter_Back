package routes

import (
	"startup_back/internal/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes ( app *fiber.App ,handlers *handler.Handlers){

	app.Use(cors.New(cors.Config{

		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
	}))
	
	api := app.Group("/api")
	auth:=api.Group("/auth")
	auth.Post("/register",handlers.User.Register)

}
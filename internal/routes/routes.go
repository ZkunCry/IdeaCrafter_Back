package routes

import (
	"startup_back/internal/handler"
	"startup_back/internal/middleware"
	"startup_back/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes ( app *fiber.App ,handlers *handler.Handlers,services *service.Services){

	app.Use(cors.New(cors.Config{

		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	api := app.Group("/api")
	auth:=api.Group("/auth")
	auth.Post("/signup", handlers.Auth.SignUp)
	auth.Post("/signin", handlers.Auth.SignIn)

	startup:=api.Group("/startup")
	startup.Post("/", middleware.RequireAuth(services.Token), handlers.Startup.CreateStartup)
	startup.Get("/list", handlers.Startup.GetListStartups)


	vacancy := api.Group("/vacancy")
	vacancy.Post("/", middleware.RequireAuth(services.Token), handlers.Vacancy.CreateVacancy)
	vacancy.Get("/:id", middleware.RequireAuth(services.Token), handlers.Vacancy.GetVacancyByID)
	vacancy.Get("/startup/:id", middleware.RequireAuth(services.Token), handlers.Vacancy.GetVacanciesByStartup)
	vacancy.Put("/:id", middleware.RequireAuth(services.Token), handlers.Vacancy.UpdateVacancy)

}
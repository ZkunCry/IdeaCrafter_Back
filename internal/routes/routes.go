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
	auth.Get("/me", handlers.Auth.IdentityMe)
	auth.Post("/refresh", handlers.Auth.Refresh)
	auth.Post("/logout", handlers.Auth.LogOut)

	startup:=api.Group("/startup")
	startup.Post("/", middleware.RequireAuth(services.Token), handlers.Startup.CreateStartup)
	startup.Get("/my-startups",middleware.RequireAuth(services.Token), handlers.Startup.GetUserStartups)
	startup.Get("/list", handlers.Startup.GetListStartups)
	startup.Get("/:id", handlers.Startup.GetStartupByID)
	startup.Post("/:id/categories", middleware.RequireAuth(services.Token), handlers.Startup.AddCategories)

	vacancy := api.Group("/vacancy")
	vacancy.Post("/", middleware.RequireAuth(services.Token), handlers.Vacancy.CreateVacancy)
	vacancy.Get("/:id",  handlers.Vacancy.GetVacancyByID)
	vacancy.Get("/startup/:id",  handlers.Vacancy.GetVacanciesByStartup)
	vacancy.Put("/:id", middleware.RequireAuth(services.Token), handlers.Vacancy.UpdateVacancy)
	vacancy.Delete("/:id", middleware.RequireAuth(services.Token), handlers.Vacancy.DeleteVacancy)


	role := api.Group("/role")
	role.Post("/", middleware.RequireAuth(services.Token), handlers.Role.CreateRole)

	application := api.Group("/application")
	application.Post("/", middleware.RequireAuth(services.Token), handlers.Application.CreateApplication)
	application.Get("/:id", handlers.Application.GetApplicationByID)
	application.Put("/:id", middleware.RequireAuth(services.Token), handlers.Application.UpdateApplication)
	application.Put("/status/:id", middleware.RequireAuth(services.Token), handlers.Application.UpdateApplicationStatus) 

	stage := api.Group("/stage")
	stage.Post("/", middleware.RequireAuth(services.Token), handlers.Stage.CreateStage)
	stage.Get("/", handlers.Stage.GetList)

	category := api.Group("/category")
	category.Post("/", middleware.RequireAuth(services.Token), handlers.Category.CreateCategory)
	category.Get("/list", handlers.Category.GetAllCategories)	
}

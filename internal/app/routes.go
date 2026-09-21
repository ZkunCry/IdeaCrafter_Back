package app

import (
	"startup_back/internal/auth"
	"startup_back/internal/platform/config"
	"startup_back/internal/platform/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App, handlers *Handlers, tokenService auth.TokenService, cfg *config.AppConfig) {

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.CORSOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	api := app.Group("/api")
	authGroup := api.Group("/auth")

	authGroup.Post("/signup", handlers.Auth.SignUp)
	authGroup.Post("/signin", handlers.Auth.SignIn)
	authGroup.Get("/me", handlers.Auth.IdentityMe)
	authGroup.Post("/refresh", handlers.Auth.Refresh)
	authGroup.Post("/logout", handlers.Auth.LogOut)

	startupGroup := api.Group("/startup")
	startupGroup.Post("/", middleware.RequireAuth(tokenService), handlers.Startup.CreateStartup)
	startupGroup.Get("/my-startups", middleware.RequireAuth(tokenService), handlers.Startup.GetUserStartups)
	startupGroup.Get("/list", handlers.Startup.GetListStartups)
	startupGroup.Get("/:id", handlers.Startup.GetStartupByID)
	startupGroup.Put("/:id", middleware.RequireAuth(tokenService), handlers.Startup.UpdateStartup)
	startupGroup.Post("/:id/categories", middleware.RequireAuth(tokenService), handlers.Startup.AddCategories)

	vacancyGroup := api.Group("/vacancy")
	vacancyGroup.Post("/", middleware.RequireAuth(tokenService), handlers.Vacancy.CreateVacancy)
	vacancyGroup.Get("/:id", handlers.Vacancy.GetVacancyByID)
	vacancyGroup.Get("/startup/:id", handlers.Vacancy.GetVacanciesByStartup)
	vacancyGroup.Put("/:id", middleware.RequireAuth(tokenService), handlers.Vacancy.UpdateVacancy)
	vacancyGroup.Delete("/:id", middleware.RequireAuth(tokenService), handlers.Vacancy.DeleteVacancy)

	roleGroup := api.Group("/role")
	roleGroup.Post("/", middleware.RequireAuth(tokenService), handlers.Role.CreateRole)
	roleGroup.Get("/", handlers.Role.GetRoles)

	applicationGroup := api.Group("/application")
	applicationGroup.Post("/", middleware.RequireAuth(tokenService), handlers.Application.CreateApplication)
	applicationGroup.Get("/my", middleware.RequireAuth(tokenService), handlers.Application.GetMyApplications)
	applicationGroup.Get("/startup/:id", middleware.RequireAuth(tokenService), handlers.Application.GetStartupApplications)
	applicationGroup.Put("/status/:id", middleware.RequireAuth(tokenService), handlers.Application.UpdateApplicationStatus)
	applicationGroup.Get("/:id", handlers.Application.GetApplicationByID)
	applicationGroup.Put("/:id", middleware.RequireAuth(tokenService), handlers.Application.UpdateApplication)
	applicationGroup.Delete("/:id", middleware.RequireAuth(tokenService), handlers.Application.DeleteApplication)

	stageGroup := api.Group("/stage")
	stageGroup.Post("/", middleware.RequireAuth(tokenService), handlers.Stage.CreateStage)
	stageGroup.Get("/", handlers.Stage.GetList)

	favoriteGroup := api.Group("/favorite")
	favoriteGroup.Get("/my", middleware.RequireAuth(tokenService), handlers.Favorite.GetMyFavorites)
	favoriteGroup.Get("/ids", middleware.RequireAuth(tokenService), handlers.Favorite.GetMyFavoriteIDs)
	favoriteGroup.Get("/startup/:id/count", handlers.Favorite.GetStartupFavoritesCount)
	favoriteGroup.Post("/startup/:id", middleware.RequireAuth(tokenService), handlers.Favorite.AddFavorite)
	favoriteGroup.Delete("/startup/:id", middleware.RequireAuth(tokenService), handlers.Favorite.RemoveFavorite)

	categoryGroup := api.Group("/category")
	categoryGroup.Post("/", middleware.RequireAuth(tokenService), handlers.Category.CreateCategory)
	categoryGroup.Get("/list", handlers.Category.GetAllCategories)
	categoryGroup.Get("/:id", handlers.Category.GetCategory)
}

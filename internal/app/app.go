package app

import (
	"context"
	"time"

	"startup_back/internal/application"
	"startup_back/internal/auth"
	"startup_back/internal/category"
	"startup_back/internal/entity"
	"startup_back/internal/favorite"
	"startup_back/internal/platform/config"
	platformdb "startup_back/internal/platform/db"
	"startup_back/internal/role"
	"startup_back/internal/stage"
	"startup_back/internal/startup"
	"startup_back/internal/user"
	"startup_back/internal/vacancy"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

type App struct {
	Fiber *fiber.App
	DB    *gorm.DB
}

func (a *App) Shutdown(ctx context.Context) error {
	shutdownErr := a.Fiber.ShutdownWithContext(ctx)

	if sqlDB, err := a.DB.DB(); err == nil {
		if closeErr := sqlDB.Close(); closeErr != nil && shutdownErr == nil {
			shutdownErr = closeErr
		}
	}

	return shutdownErr
}

type Handlers struct {
	Auth        *auth.Handler
	Startup     *startup.Handler
	Vacancy     *vacancy.Handler
	Role        *role.Handler
	Application *application.Handler
	Stage       *stage.Handler
	Category    *category.Handler
	Favorite    *favorite.Handler
}

func New(cfg *config.AppConfig) (*App, error) {
	db, err := platformdb.Open(cfg)
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Startup{},
		&entity.Category{},
		&entity.Favorite{},
		&entity.StartupFile{},
		&entity.Role{},
		&entity.Vacancy{},
		&entity.Application{},
		&entity.Stage{},
	); err != nil {
		return nil, err
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	tokenService := auth.NewTokenService(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret)
	passwordService := auth.NewPasswordService()
	authService := auth.NewService(userService, passwordService, tokenService)

	startupRepository := startup.NewRepository(db)
	startupService := startup.NewService(startupRepository)

	vacancyRepository := vacancy.NewRepository(db)
	vacancyService := vacancy.NewService(vacancyRepository)

	roleRepository := role.NewRepository(db)
	roleService := role.NewService(roleRepository)

	applicationRepository := application.NewRepository(db)
	applicationService := application.NewService(applicationRepository)

	stageRepository := stage.NewRepository(db)
	stageService := stage.NewService(stageRepository)

	categoryRepository := category.NewRepository(db)
	categoryService := category.NewService(categoryRepository)

	favoriteRepository := favorite.NewRepository(db)
	favoriteService := favorite.NewService(favoriteRepository)

	handlers := &Handlers{
		Auth:        auth.NewHandler(authService),
		Startup:     startup.NewHandler(startupService, cfg.S3Client, cfg.S3.Bucket),
		Vacancy:     vacancy.NewHandler(vacancyService),
		Role:        role.NewHandler(roleService),
		Application: application.NewHandler(applicationService),
		Stage:       stage.NewHandler(stageService),
		Category:    category.NewHandler(categoryService),
		Favorite:    favorite.NewHandler(favoriteService),
	}

	app := fiber.New(fiber.Config{
		AppName:               "startup_back",
		ReadTimeout:           15 * time.Second,
		WriteTimeout:          30 * time.Second,
		IdleTimeout:           60 * time.Second,
		BodyLimit:             32 * 1024 * 1024,
		DisableStartupMessage: cfg.IsProduction(),
		ProxyHeader: fiber.HeaderXForwardedFor,
	})

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(fiberlogger.New())

	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe:     func(c *fiber.Ctx) bool { return true },
		LivenessEndpoint:  "/health",
		ReadinessProbe:    func(c *fiber.Ctx) bool { return pingDB(c.Context(), db) },
		ReadinessEndpoint: "/ready",
	}))

	if !cfg.IsProduction() {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	SetupRoutes(app, handlers, tokenService, cfg)

	return &App{Fiber: app, DB: db}, nil
}

func pingDB(ctx context.Context, db *gorm.DB) bool {
	sqlDB, err := db.DB()
	if err != nil {
		return false
	}

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return sqlDB.PingContext(pingCtx) == nil
}

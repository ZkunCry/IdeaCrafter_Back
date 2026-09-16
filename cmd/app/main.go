package main

import (
	"fmt"
	application "startup_back/internal/app"
	"startup_back/internal/platform/config"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.LoadConfig()

	logrus.SetFormatter(&logrus.JSONFormatter{})
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	app, err := application.New(cfg)
	if err != nil {
		logrus.Fatalf("failed to build application: %v", err)
	}

	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logrus.Infof("Starting server on %s", address)
	if err := app.Listen(address); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

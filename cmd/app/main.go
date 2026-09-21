package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	application "startup_back/internal/app"
	"startup_back/internal/platform/config"

	"github.com/sirupsen/logrus"
)

var version = "dev"

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
	logrus.SetOutput(os.Stdout)

	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	if !cfg.IsProduction() {
		logrus.SetLevel(logrus.DebugLevel)
	}

	app, err := application.New(cfg)
	if err != nil {
		logrus.Fatalf("failed to build application: %v", err)
	}

	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	serverErr := make(chan error, 1)
	go func() {
		logrus.WithFields(logrus.Fields{
			"address": address,
			"env":     cfg.Env,
			"version": version,
		}).Info("starting server")

		if err := app.Fiber.Listen(address); err != nil {
			serverErr <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		logrus.Fatalf("failed to start server: %v", err)
	case sig := <-stop:
		logrus.Infof("received %s, shutting down", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		logrus.Errorf("graceful shutdown failed: %v", err)
		os.Exit(1)
	}

	logrus.Info("server stopped")
}

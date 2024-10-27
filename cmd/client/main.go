package main

import (
	"github.com/SversusN/keeper/internal/client/app"
	"github.com/SversusN/keeper/internal/client/config"
	"github.com/SversusN/keeper/pkg/logger"
	"log"
)

const (
	configPath = "./config/client/appsettings.json"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	settings, err := config.Initialize(configPath)
	if err != nil {
		return err
	}
	l, err := logger.Initialize(settings.LogLevel)
	if err != nil {
		return err
	}

	client, err := app.NewClient(l, settings, buildVersion, buildDate)
	if err != nil {
		return err
	}

	return client.Start()
}

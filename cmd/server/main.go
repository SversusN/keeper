// Keeper - приложение, реализующее gRPC сервер для хранения зашифрованной информации
package main

import (
	"github.com/SversusN/keeper/internal/client/utils"
	"log"

	"github.com/SversusN/keeper/internal/server/config"
	"github.com/SversusN/keeper/internal/server/handlers"
	"github.com/SversusN/keeper/internal/server/storage"
	"github.com/SversusN/keeper/pkg/logger"
)

const (
	configPath = "./config/server/appsettings.json"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	settings, err := config.Initialize(configPath)
	if err != nil {
		settings, err = config.Initialize(utils.SettingsParentPath())
		if err != nil {
			return err
		}
	}
	l, err := logger.Initialize(settings.LogLevel)
	if err != nil {
		return err
	}

	rep, err := storage.NewStorage(settings.DatabaseURL, l)
	if err != nil {
		l.Log.Panicf("db error: %v", err)
		return err
	}
	defer func() {
		if err = rep.Close(); err != nil {
			l.Log.Errorf("db close error: %v", err)
		}
	}()

	serverGRPC := handlers.NewServer(rep, settings, l)

	return serverGRPC.Start()
}

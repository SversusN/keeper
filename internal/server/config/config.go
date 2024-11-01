package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/caarlos0/env"
)

// Config – объект конфигурации сервера.
type Config struct {
	// 	DatabaseURL  – dsn для подключения к БД.
	DatabaseURL string `json:"database_url"`
	// Host – адрес сервера.
	Host string `json:"host"`
	// LogLevel – уровень логгирования.
	LogLevel string `json:"log_level"`
	// SecretKey – ключ шифрования.
	SecretKey string `json:"secret_key"`
}

// Initialize – функция инициализации конфига.
func Initialize(configPath string) (*Config, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("open file error: %w", err)
	}
	var c = &Config{}
	err = json.NewDecoder(configFile).Decode(c)
	if err != nil {
		return nil, fmt.Errorf("parse JSON error: %w", err)
	}
	err = env.Parse(c)
	return c, err
}

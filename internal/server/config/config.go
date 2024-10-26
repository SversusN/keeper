package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/caarlos0/env"
)

// Config – объект конфигурации сервера.
type Config struct {
	// DatabaseURL – dsn для подключения к БД.
	DatabaseURL string `json:"database_url" env:"DATABASE_URL"`
	// Host – адрес сервера.
	Host string `json:"host" env:"HOST"`
	// LogLevel – уровень логгирования.
	LogLevel string `json:"log_level" env:"LOG_LEVEL"`
	// SecretKey – ключ шифрования.
	SecretKey string `json:"secret_key" env:"SECRET_KEY"`
}

// Initialize – функция инициализации конфига.
func Initialize(configPath string) (*Config, error) {
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("open file error: %w", err)
	}
	var c = &Config{}
	err = json.Unmarshal(configFile, c)
	if err != nil {
		return nil, fmt.Errorf("parse JSON error: %w", err)
	}
	err = env.Parse(c)
	return c, err
}

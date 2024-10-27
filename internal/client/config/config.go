package config

import (
	"encoding/json"
	"os"

	"github.com/caarlos0/env"
)

// Config – объект конфигурации клиента.
type Config struct {
	// Host – адрес сервера.
	Host string `json:"host" env:"HOST"`
	// LogLevel – уровень логирования.
	LogLevel string `json:"log_level" env:"LOG_LEVEL"`
	// ConnectionTimeout – таймаут запроса к серверу.
	ConnectionTimeout int64 `json:"connection_timeout" env:"CONNECTION_TIMEOUT"`
	// ChanSize – размер канала для синхронизации данных.
	ChanSize int64 `json:"chan_size" env:"CHAN_SIZE"`
	//Секретное слово для клиента
	PassPhrase string `json:"passphrase" env:"PASS_PHRASE"`
}

// Initialize – функция инициализации конфига.
func Initialize(configPath string) (*Config, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	var c = &Config{}
	err = json.NewDecoder(configFile).Decode(c)
	if err != nil {
		return nil, err
	}
	err = env.Parse(c)
	return c, err
}

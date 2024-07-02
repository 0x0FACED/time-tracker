package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Username string
	Name     string
	Host     string
	Port     string
	Pass     string
	Driver   string
}

type ServerConfig struct {
	Host string
	Port string
}

type Config struct {
	DatabaseConfig
	ServerConfig
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		DatabaseConfig{
			Username: os.Getenv("DB_USER"),
			Pass:     os.Getenv("DB_PASS"),
			Name:     os.Getenv("DB_NAME"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Driver:   os.Getenv("DB_DRIVER"),
		},
		ServerConfig{
			Host: os.Getenv("S_HOST"),
			Port: os.Getenv("S_PORT"),
		},
	}, nil
}

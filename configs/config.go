package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	DBUsername string
	DBName     string
	DBHost     string
	DBPort     string
	DBPass     string
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
			DBUsername: os.Getenv("DB_USER"),
			DBPass:     os.Getenv("DB_PASS"),
			DBName:     os.Getenv("DB_NAME"),
			DBHost:     os.Getenv("DB_HOST"),
			DBPort:     os.Getenv("DB_PORT"),
		},
		ServerConfig{
			Host: os.Getenv("S_HOST"),
			Port: os.Getenv("S_PORT"),
		},
	}, nil
}

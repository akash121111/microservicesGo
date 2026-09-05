package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	POSTGRES_HOST     string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string
	POSTGRES_PORT     string
	REDIS_HOST        string
	REDIS_PORT        string
	RABBITMQ_HOST     string
	RABBITMQ_PORT     string
	RABBITMQ_USER     string
	RABBITMQ_PASSWORD string
	GRPCPort          string
}

func LoadConfig() (*Config, error) {
	// .env is optional.
	_ = godotenv.Load()

	cfg := &Config{
		Port:              os.Getenv("PORT"),
		POSTGRES_HOST:     os.Getenv("POSTGRES_HOST"),
		POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		POSTGRES_DB:       os.Getenv("POSTGRES_DB"),
		POSTGRES_PORT:     os.Getenv("POSTGRES_PORT"),
		REDIS_HOST:        os.Getenv("REDIS_HOST"),
		REDIS_PORT:        os.Getenv("REDIS_PORT"),
		RABBITMQ_HOST:     os.Getenv("RABBITMQ_HOST"),
		RABBITMQ_PORT:     os.Getenv("RABBITMQ_PORT"),
		RABBITMQ_USER:     os.Getenv("RABBITMQ_USER"),
		RABBITMQ_PASSWORD: os.Getenv("RABBITMQ_PASSWORD"),
		GRPCPort:          os.Getenv("GRPC_PORT"),
	}

	if cfg.Port == "" {
		return nil, fmt.Errorf("PORT is required")
	}

	if cfg.POSTGRES_HOST == "" {
		return nil, fmt.Errorf("POSTGRES_HOST is required")
	}

	if cfg.POSTGRES_USER == "" {
		return nil, fmt.Errorf("POSTGRES_USER is required")
	}

	if cfg.POSTGRES_PASSWORD == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD is required")
	}

	if cfg.POSTGRES_DB == "" {
		return nil, fmt.Errorf("POSTGRES_DB is required")
	}

	if cfg.POSTGRES_PORT == "" {
		return nil, fmt.Errorf("POSTGRES_PORT is required")
	}
	if cfg.REDIS_HOST == "" {
		return nil, fmt.Errorf("REDIS_HOST is required")
	}
	if cfg.REDIS_PORT == "" {
		return nil, fmt.Errorf("REDIS_PORT is required")
	}

	return cfg, nil
}

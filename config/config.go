package config

import (
	"log/slog"

	env "github.com/caarlos0/env/v6"
)

type environment struct {
	BrightDataHost     string `env:"BRIGHT_DATA_HOST,required"`
	BrightDataPort     int    `env:"BRIGHT_DATA_PORT,required"`
	BrightDataUsername string `env:"BRIGHT_DATA_USERNAME,required"`
	BrightDataPassword string `env:"BRIGHT_DATA_PASSWORD,required"`
	RabbitMQUser       string `env:"RABBITMQ_USER,required"`
	RabbitMQPassword   string `env:"RABBITMQ_PASSWORD,required"`
	RabbitMQHost       string `env:"RABBITMQ_HOST,required"`
	RabbitMQPort       string `env:"RABBITMQ_PORT,required"`
}

func NewConfig() (Config, error) {
	slog.Info("Loading environment...")
	environment := environment{}
	if err := env.Parse(&environment); err != nil {
		return Config{}, err
	}

	slog.Info("Environment loaded successfully!")

	cfg := Config{
		BrightDataConfig: brightDataConfig{
			Host:     environment.BrightDataHost,
			Port:     environment.BrightDataPort,
			Username: environment.BrightDataUsername,
			Password: environment.BrightDataPassword,
		},
		RabbitMQConfig: rabbitMQConfig{
			Host:     environment.RabbitMQHost,
			Port:     environment.RabbitMQPort,
			Username: environment.RabbitMQUser,
			Password: environment.RabbitMQPassword,
		},
	}

	return cfg, nil
}

type Config struct {
	BrightDataConfig brightDataConfig
	RabbitMQConfig   rabbitMQConfig
}

type brightDataConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type rabbitMQConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

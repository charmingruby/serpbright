package config

import (
	"log/slog"

	env "github.com/caarlos0/env/v6"
)

type environment struct {
	BrightDataHost       string `env:"BRIGHT_DATA_HOST,required"`
	BrightDataPort       int    `env:"BRIGHT_DATA_PORT,required"`
	BrightDataCustomerID string `env:"BRIGHT_DATA_CUSTOMER_ID,required"`
	BrightDataZone       string `env:"BRIGHT_DATA_ZONE,required"`
	BrightDataPassword   string `env:"BRIGHT_DATA_PASSWORD,required"`
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
			Host:       environment.BrightDataHost,
			Port:       environment.BrightDataPort,
			CustomerID: environment.BrightDataCustomerID,
			Zone:       environment.BrightDataZone,
			Password:   environment.BrightDataPassword,
		},
	}

	return cfg, nil
}

type Config struct {
	BrightDataConfig brightDataConfig
}

type brightDataConfig struct {
	Host       string
	Port       int
	CustomerID string
	Zone       string
	Username   string
	Password   string
}

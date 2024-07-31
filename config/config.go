package config

import (
	"log/slog"

	env "github.com/caarlos0/env/v6"
)

type environment struct {
	BrightDataAPIKey string `env:"BRIGHT_DATA_API_KEY,required"`
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
			APIKey: environment.BrightDataAPIKey,
		},
	}

	return cfg, nil
}

type Config struct {
	BrightDataConfig brightDataConfig
}

type brightDataConfig struct {
	APIKey string
}

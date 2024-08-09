package config

import (
	"log/slog"

	env "github.com/caarlos0/env/v6"
)

type environment struct {
	SearchResultIncludeHTML bool   `env:"SEARCH_RESULT_INCLUDE_HTML,required"`
	BrightDataHost          string `env:"BRIGHT_DATA_HOST,required"`
	BrightDataPort          int    `env:"BRIGHT_DATA_PORT,required"`
	BrightDataUsername      string `env:"BRIGHT_DATA_USERNAME,required"`
	BrightDataPassword      string `env:"BRIGHT_DATA_PASSWORD,required"`
	RabbitMQURI             string `env:"RABBITMQ_URI,required"`
	MongoURI                string `env:"MONGO_URI,required"`
	MongoDatabase           string `env:"MONGO_DATABASE,required"`
	DebugMode               bool   `env:"DEBUG_MODE,required"`
}

func NewConfig() (Config, error) {
	slog.Info("ENVIRONMENT: " + "Loading environment variables...")
	environment := environment{}
	if err := env.Parse(&environment); err != nil {
		return Config{}, err
	}

	slog.Info("ENVIRONMENT: Environment variables loaded successfully!")

	cfg := Config{
		DebugMode: environment.DebugMode,
		BrightDataConfig: brightDataConfig{
			Host:     environment.BrightDataHost,
			Port:     environment.BrightDataPort,
			Username: environment.BrightDataUsername,
			Password: environment.BrightDataPassword,
		},
		RabbitMQConfig: rabbitMQConfig{
			URI: environment.RabbitMQURI,
		},
		MongoConfig: mongoConfig{
			URI:          environment.MongoURI,
			DatabaseName: environment.MongoDatabase,
		},
		SearchConfig: searchConfig{
			IncludeHTML: environment.SearchResultIncludeHTML,
		},
	}

	return cfg, nil
}

type Config struct {
	DebugMode        bool
	BrightDataConfig brightDataConfig
	RabbitMQConfig   rabbitMQConfig
	MongoConfig      mongoConfig
	SearchConfig     searchConfig
}

type brightDataConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type rabbitMQConfig struct {
	URI string
}

type mongoConfig struct {
	URI          string
	DatabaseName string
}

type searchConfig struct {
	IncludeHTML bool
}

package config

import (
	"log/slog"

	env "github.com/caarlos0/env/v6"
)

type environment struct {
	BrightDataHost                   string `env:"BRIGHT_DATA_HOST,required"`
	BrightDataPort                   int    `env:"BRIGHT_DATA_PORT,required"`
	BrightDataUsername               string `env:"BRIGHT_DATA_USERNAME,required"`
	BrightDataPassword               string `env:"BRIGHT_DATA_PASSWORD,required"`
	RabbitMQURI                      string `env:"RABBITMQ_URI,required"`
	MongoURI                         string `env:"MONGO_URI,required"`
	MongoDatabase                    string `env:"MONGO_DATABASE,required"`
	SearchIncludeHTML                bool   `env:"SEARCH_INCLUDE_HTML,required"`
	SearchSkipRedirectAll            bool   `env:"SEARCH_SKIP_REDIRECT_ALL,required"`
	SearchSkipRedirectCampaigns      string `env:"SEARCH_SKIP_REDIRECT_CAMPAIGNS,required"`
	SearchSkipCustomerDomainRedirect bool   `env:"SEARCH_SKIP_CUSTOMER_DOMAIN_REDIRECT,required"`
	SearchRedirectTimeout            int    `env:"SEARCH_REDIRECT_TIMEOUT,required"`
	SearchConcatFirstURLDomain       string `env:"SEARCH_CONCAT_FIRST_URL_DOMAIN,required"`
	SearchConcatDomainLastURL        string `env:"SEARCH_CONCAT_DOMAIN_LAST_URL,required"`
	DebugMode                        bool   `env:"DEBUG_MODE,required"`
	ConcurrencyThreads               int    `env:"CONCURRENCY_THREADS,required"`
}

func NewConfig() (Config, error) {
	slog.Info("ENVIRONMENT: " + "Loading environment variables...")
	environment := environment{}
	if err := env.Parse(&environment); err != nil {
		return Config{}, err
	}

	slog.Info("ENVIRONMENT: Environment variables loaded successfully!")

	cfg := Config{
		DebugMode:          environment.DebugMode,
		ConcurrencyThreads: environment.ConcurrencyThreads,
		BrightDataConfig: BrightDataConfig{
			Host:     environment.BrightDataHost,
			Port:     environment.BrightDataPort,
			Username: environment.BrightDataUsername,
			Password: environment.BrightDataPassword,
		},
		RabbitMQConfig: RabbitMQConfig{
			URI: environment.RabbitMQURI,
		},
		MongoConfig: MongoConfig{
			URI:          environment.MongoURI,
			DatabaseName: environment.MongoDatabase,
		},
		SearchConfig: SearchConfig{
			IncludeHTML:                environment.SearchIncludeHTML,
			SkipRedirectAll:            environment.SearchSkipRedirectAll,
			SkipRedirectCampaigns:      environment.SearchSkipRedirectCampaigns,
			SkipCustomerDomainRedirect: environment.SearchSkipCustomerDomainRedirect,
			RedirectTimeout:            environment.SearchRedirectTimeout,
			ConcatFirstDomainURL:       environment.SearchConcatFirstURLDomain,
			ConcatDomainLastURL:        environment.SearchConcatDomainLastURL,
		},
	}

	return cfg, nil
}

type Config struct {
	DebugMode          bool
	ConcurrencyThreads int
	BrightDataConfig   BrightDataConfig
	RabbitMQConfig     RabbitMQConfig
	MongoConfig        MongoConfig
	SearchConfig       SearchConfig
}

type BrightDataConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type RabbitMQConfig struct {
	URI string
}

type MongoConfig struct {
	URI          string
	DatabaseName string
}

type SearchConfig struct {
	IncludeHTML                bool
	SkipRedirectAll            bool
	SkipRedirectCampaigns      string
	SkipCustomerDomainRedirect bool
	RedirectTimeout            int
	ConcatDomainLastURL        string
	ConcatFirstDomainURL       string
}

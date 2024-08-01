package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/charmingruby/serpright/config"
	"github.com/charmingruby/serpright/internal/scrapper"
	"github.com/charmingruby/serpright/internal/scrapper/domain/usecase"
	"github.com/charmingruby/serpright/internal/scrapper/serp/service/brightdata"
	"github.com/charmingruby/serpright/test/inmemory"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("CONFIGURATION: .env file not found")
	}

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error(fmt.Sprintf("CONFIGURATION: %s", err.Error()))
		os.Exit(1)
	}

	serp := brightdata.NewBrightData(cfg)
	campaingTaskRepo := inmemory.NewInMemoryCampaignTaskRepository()
	svc := scrapper.NewService(serp, &campaingTaskRepo)

	runBrightDataActions(svc, serp)
}

func runBrightDataActions(svc usecase.ScrapperUseCase, brightData *brightdata.BrightData) {
	slog.Info("Running BrightData actions...")

	_, err := brightData.ExecSearch(svc)
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err.Error()))
	}
}

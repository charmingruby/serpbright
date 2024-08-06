package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/charmingruby/serpright/config"
	rabbitMQPubSub "github.com/charmingruby/serpright/internal/common/queue/rabbitmq"
	"github.com/charmingruby/serpright/internal/scrapper"
	"github.com/charmingruby/serpright/internal/scrapper/domain/event"
	"github.com/charmingruby/serpright/internal/scrapper/domain/usecase"
	"github.com/charmingruby/serpright/internal/scrapper/infra/queue"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/brightdata"
	"github.com/charmingruby/serpright/pkg/rabbitmq"

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

	ch, close := rabbitmq.NewRabbitMQConnection(&cfg)
	defer func() {
		close()
		ch.Close()
	}()

	pubsub := rabbitMQPubSub.NewRabbitMQPubSub(ch)

	serp := brightdata.NewBrightData(cfg)
	svc := scrapper.NewService(serp)

	processCampaingTaskEventHandler := queue.NewCampaignTaskProcessHandler(svc)
	go pubsub.Subscribe(event.ProcessCampaignTask, processCampaingTaskEventHandler.Handle)

	runBrightDataActions(svc, serp, cfg.DebugMode)
}

func runBrightDataActions(
	svc usecase.ScrapperUseCase,
	brightData *brightdata.BrightData,
	debug bool) {
	slog.Info("Running BrightData actions...")

	op, err := brightData.ExecSearch(svc, debug)
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err.Error()))
	}

	// Example processed data
	fmt.Printf("RequestID: %s\n", op.SearchResult.RequestID)
	fmt.Printf("SearchType: %s\n", op.SearchResult.SearchType)
}

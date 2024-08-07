package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmingruby/serpright/config"
	rabbitMQPubSub "github.com/charmingruby/serpright/internal/common/queue/rabbitmq"
	"github.com/charmingruby/serpright/internal/scrapper"
	"github.com/charmingruby/serpright/internal/scrapper/domain/event"
	"github.com/charmingruby/serpright/internal/scrapper/infra/database/mongo_repository"
	"github.com/charmingruby/serpright/internal/scrapper/infra/queue"
	"github.com/charmingruby/serpright/internal/scrapper/infra/serp/brightdata"
	mongodb "github.com/charmingruby/serpright/pkg/mongo"
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

	ch, close := rabbitmq.NewRabbitMQConnection(cfg.RabbitMQConfig.URI)
	defer func() {
		close()
		ch.Close()
	}()

	db, err := mongodb.NewMongoConnection(cfg.MongoConfig.URI, cfg.MongoConfig.DatabaseName)
	if err != nil {
		slog.Error(fmt.Sprintf("MONGO CONNECTION: %s", err.Error()))
		os.Exit(1)
	}

	pubsub := rabbitMQPubSub.NewRabbitMQPubSub(ch)
	searchResultRepo := mongo_repository.NewSearchResultMongoRepository(db)

	serp := brightdata.NewBrightData(cfg)
	svc := scrapper.NewService(serp, &searchResultRepo)

	processCampaingTaskEventHandler := queue.NewCampaignTaskProcessHandler(svc, cfg.DebugMode)

	go pubsub.Subscribe(event.ProcessCampaignTask, processCampaingTaskEventHandler.Handle)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	slog.Info("Terminating gracefully")
}

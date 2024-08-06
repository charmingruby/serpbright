package rabbitmq

import (
	"log"
	"log/slog"

	"github.com/charmingruby/serpright/config"
	"github.com/charmingruby/serpright/internal/scrapper/domain/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQConnection(cfg *config.Config) (*amqp.Channel, func() error) {
	addr := cfg.RabbitMQConfig.URI

	conn, err := amqp.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	err = ch.ExchangeDeclare(event.ProcessCampaignTask, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Connected successfully to RabbitMQ!")

	return ch, ch.Close
}

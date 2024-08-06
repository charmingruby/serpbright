package rabbitmq

import (
	"log"
	"log/slog"

	"github.com/charmingruby/serpright/internal/scrapper/domain/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQConnection(uri string) (*amqp.Channel, func() error) {
	conn, err := amqp.Dial(uri)
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

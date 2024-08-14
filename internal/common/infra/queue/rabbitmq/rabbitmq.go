package rabbitmq

import (
	"context"
	"log"
	"log/slog"

	"github.com/charmingruby/serpright/internal/scrapper/domain/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQPubSub(ch *amqp.Channel) RabbitMQPubSub {
	return RabbitMQPubSub{
		ch: ch,
	}
}

type RabbitMQPubSub struct {
	ch *amqp.Channel
}

func (rmq *RabbitMQPubSub) Publish(topic string, message []byte) error {
	q, err := rmq.ch.QueueDeclare(topic, true, false, false, false, nil)
	if err != nil {
		return err
	}

	if err := rmq.ch.PublishWithContext(context.Background(), "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         message,
		DeliveryMode: amqp.Persistent,
	}); err != nil {
		return err
	}

	log.Printf("Published message to %s -> %v", topic, string(message))

	return nil
}

func (rmq *RabbitMQPubSub) Subscribe(topic string, handler func([]byte)) {
	q, err := rmq.ch.QueueDeclare(event.ProcessCampaignTask, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := rmq.ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan struct{})

	slog.Info("RABBITMQ: Pooling on " + topic)

	go func() {
		for d := range msgs {
			log.Printf("Received message from %s -> %v", topic, string(d.MessageId))
			handler(d.Body)
		}
	}()

	<-forever
}

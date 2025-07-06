package utils

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func (r *RabbitMQ) DeclareQueue(name string) error {
	_, err := r.Channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (r *RabbitMQ) Publish(queueName string, body []byte, headers amqp.Table) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.Channel.PublishWithContext(ctx,
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Headers:     headers,
		},
	)
}

func (r *RabbitMQ) GetRetryCount(headers amqp.Table) int {
	if val, ok := headers["x-retry"]; ok {
		switch v := val.(type) {
		case int32:
			return int(v)
		case int64:
			return int(v)
		case int:
			return v
		}
	}
	return 0
}

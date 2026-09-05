package messaging

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQ) PublishOutboxEvent(
	ctx context.Context,
	routingKey string,
	payload []byte,
) error {

	err := r.Ch.PublishWithContext(
		ctx,
		OrderExchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp091.Persistent,
			Body:         payload,
		},
	)

	if err != nil {
		return fmt.Errorf(
			"failed to publish outbox event: %w",
			err,
		)
	}

	return nil
}

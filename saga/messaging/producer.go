package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQ) PublishReserveStock(
	ctx context.Context,
	command ReserveStockCommand,
) error {

	body, err := json.Marshal(command)
	if err != nil {
		return fmt.Errorf(
			"failed to marshal reserve stock command: %w",
			err,
		)
	}

	err = r.Ch.PublishWithContext(
		ctx,
		OrderExchange,
		ReserveStockRoutingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp091.Persistent,
			Body:         body,
		},
	)

	if err != nil {
		return fmt.Errorf(
			"failed to publish reserve stock command: %w",
			err,
		)
	}

	return nil
}

func (r *RabbitMQ) PublishPament(ctx context.Context, command PaymentProcessCommand) error {
	body, err := json.Marshal(command)

	if err != nil {
		return fmt.Errorf(
			"failed to marshal payment command: %w",
			err,
		)
	}
	err = r.Ch.PublishWithContext(ctx, OrderExchange, PaymentProcessRoutingKey, false, false, amqp091.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp091.Persistent,
		Body:         body,
	})
	if err != nil {
		return fmt.Errorf("failed to publish pament process command %w", err)
	}
	return nil
}

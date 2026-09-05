package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQ) PublishStockReservation(
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
		StockReservedRoutingKey,
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

func (r *RabbitMQ) PublishStockReservationFailed(
	ctx context.Context,
	command ReserveStockFailedEvent,
) error {

	body, err := json.Marshal(command)
	if err != nil {
		return fmt.Errorf(
			"failed to marshal reserve stock failed command: %w",
			err,
		)
	}

	err = r.Ch.PublishWithContext(
		ctx,
		OrderExchange,
		ReservationStockFailedRoutingKey,
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

func (r *RabbitMQ) PublishToDLQ(
	message amqp091.Delivery,
) error {

	return r.Ch.Publish(
		StockDLX,
		StockDLQ,
		false,
		false,
		amqp091.Publishing{
			ContentType:  message.ContentType,
			DeliveryMode: amqp091.Persistent,
			Body:         message.Body,
			Headers:      message.Headers,
		},
	)
}

package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQ) PublishPaymentSuccess(
	ctx context.Context,
	event PaymentSuccessEvent,
) error {

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf(
			"failed to marshal payment.success: %w",
			err,
		)
	}

	err = r.Ch.PublishWithContext(
		ctx,
		OrderExchange,
		PaymentSuccessRoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)

	if err != nil {
		return fmt.Errorf(
			"failed to publish payment.success: %w",
			err,
		)
	}

	return nil
}

func (r *RabbitMQ) PublishPaymentFailed(
	ctx context.Context,
	event PaymentFailedEvent,
) error {

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf(
			"failed to marshal payment.failed: %w",
			err,
		)
	}

	err = r.Ch.PublishWithContext(
		ctx,
		OrderExchange,
		PaymentFailedRoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)

	if err != nil {
		return fmt.Errorf(
			"failed to publish payment.failed: %w",
			err,
		)
	}

	return nil
}

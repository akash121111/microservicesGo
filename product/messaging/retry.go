package messaging

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func getRetryCount(
	message amqp091.Delivery,
) int {

	value, ok := message.Headers["x-retry-count"]

	if !ok {
		return 0
	}

	switch count := value.(type) {

	case int:
		return count

	case int32:
		return int(count)

	case int64:
		return int(count)

	default:
		return 0
	}
}

func (r *RabbitMQ) RetryMessage(
	message amqp091.Delivery,
) {

	retryCount := getRetryCount(message)

	// ========================================
	// Maximum retry reached
	// ========================================

	if retryCount >= MaxStockRetries {

		log.Printf(
			"maximum retries reached: retryCount=%d",
			retryCount,
		)

		err := r.PublishToDLQ(message)

		if err != nil {

			log.Printf(
				"failed to publish message to DLQ: %v",
				err,
			)

			// If DLQ publishing itself failed,
			// requeue the original message.
			message.Nack(false, true)
			return
		}

		message.Ack(false)
		return
	}

	// ========================================
	// Increase retry count
	// ========================================

	retryCount++

	headers := amqp091.Table{}

	// Preserve existing headers.
	for key, value := range message.Headers {
		headers[key] = value
	}

	headers["x-retry-count"] = int32(retryCount)

	// ========================================
	// Publish to retry exchange
	// ========================================

	err := r.Ch.Publish(
		StockRetryExchange,
		StockRetryRoutingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType:  message.ContentType,
			DeliveryMode: amqp091.Persistent,
			Body:         message.Body,
			Headers:      headers,
		},
	)

	if err != nil {

		log.Printf(
			"failed to publish retry message: %v",
			err,
		)

		// Retry publishing the retry message.
		message.Nack(false, true)
		return
	}

	log.Printf(
		"message sent to retry queue: retryCount=%d",
		retryCount,
	)

	// Original message is now safely copied
	// into retry queue.
	message.Ack(false)
}

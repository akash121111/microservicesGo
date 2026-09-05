package messaging

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

const (
	OrderExchange = "order.events"

	ReserveStockRoutingKey           = "stock.reserve"
	StockReservedRoutingKey          = "stock.reserved"
	ReservationStockFailedRoutingKey = "stock.reservation.failed"
	StockQueue                       = "product.stock.queue"

	// Retry
	StockRetryExchange   = "product.stock.retry.exchange"
	StockRetryQueue      = "product.stock.retry.queue"
	StockRetryRoutingKey = "stock.retry"

	StockDLX = "product.stock.dlx"

	StockDLQ = "product.stock.dlq"

	MaxStockRetries = 3

	StockRetryTTL = 5000
)

type RabbitMQ struct {
	Conn *amqp091.Connection
	Ch   *amqp091.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect to rabbitmq: %w",
			err,
		)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()

		return nil, fmt.Errorf(
			"failed to create rabbitmq channel: %w",
			err,
		)
	}

	rabbit := &RabbitMQ{
		Conn: conn,
		Ch:   ch,
	}

	if err := rabbit.setup(); err != nil {
		rabbit.Close()
		return nil, err
	}

	return rabbit, nil
}

func (r *RabbitMQ) setup() error {

	// Main exchange
	err := r.Ch.ExchangeDeclare(
		OrderExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}
	err = r.Ch.ExchangeDeclare(
		StockRetryExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to declare stock retry exchange: %w",
			err,
		)
	}

	// stock Dead Letter Exchange
	err = r.Ch.ExchangeDeclare(
		StockDLX,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare stock DLX: %w", err)
	}

	// stock Dead Letter Queue
	_, err = r.Ch.QueueDeclare(
		StockDLQ,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare stock DLQ: %w", err)
	}

	// Bind stock DLQ
	err = r.Ch.QueueBind(
		StockDLQ,
		StockDLQ,
		StockDLX,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind stock DLQ: %w", err)
	}

	// stock main queue
	_, err = r.Ch.QueueDeclare(
		StockQueue,
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-dead-letter-exchange":    StockDLX,
			"x-dead-letter-routing-key": StockDLQ,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to declare stock queue: %w", err)
	}

	_, err = r.Ch.QueueDeclare(
		StockRetryQueue,
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-message-ttl":             StockRetryTTL,
			"x-dead-letter-exchange":    OrderExchange,
			"x-dead-letter-routing-key": ReserveStockRoutingKey,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"failed to declare stock retry queue: %w",
			err,
		)
	}

	err = r.Ch.QueueBind(
		StockRetryQueue,
		StockRetryRoutingKey,
		StockRetryExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to bind stock retry queue: %w",
			err,
		)
	}

	// stock consumes order.created
	err = r.Ch.QueueBind(
		StockQueue,
		ReserveStockRoutingKey,
		OrderExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind stock queue: %w", err)
	}

	// Prefetch
	err = r.Ch.Qos(
		10,
		0,
		false,
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	return nil
}

func (r *RabbitMQ) Close() {
	if r.Ch != nil {
		_ = r.Ch.Close()
	}

	if r.Conn != nil {
		_ = r.Conn.Close()
	}
}

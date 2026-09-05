package messaging

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

const (
	RabbitMQURL = "amqp://guest:guest@localhost:5672/"

	OrderExchange = "order.events"

	PaymentProcessRoutingKey = "payment.process"
	PaymentSuccessRoutingKey = "payment.success"
	PaymentFailedRoutingKey  = "payment.failed"

	PaymentQueue = "payment.queue"

	PaymentDLX = "payment.dlx"
	PaymentDLQ = "payment.dlq"
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

	// ========================================
	// Main Exchange
	// ========================================

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
		return fmt.Errorf(
			"failed to declare exchange: %w",
			err,
		)
	}

	// ========================================
	// Payment DLX
	// ========================================

	err = r.Ch.ExchangeDeclare(
		PaymentDLX,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to declare payment DLX: %w",
			err,
		)
	}

	// ========================================
	// Payment DLQ
	// ========================================

	_, err = r.Ch.QueueDeclare(
		PaymentDLQ,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to declare payment DLQ: %w",
			err,
		)
	}

	err = r.Ch.QueueBind(
		PaymentDLQ,
		PaymentDLQ,
		PaymentDLX,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to bind payment DLQ: %w",
			err,
		)
	}

	// ========================================
	// Payment Queue
	// ========================================

	_, err = r.Ch.QueueDeclare(
		PaymentQueue,
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-dead-letter-exchange":    PaymentDLX,
			"x-dead-letter-routing-key": PaymentDLQ,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"failed to declare payment queue: %w",
			err,
		)
	}

	// ========================================
	// Bind payment.process
	// ========================================

	err = r.Ch.QueueBind(
		PaymentQueue,
		PaymentProcessRoutingKey,
		OrderExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to bind payment queue: %w",
			err,
		)
	}

	// ========================================
	// Prefetch
	// ========================================

	err = r.Ch.Qos(
		10,
		0,
		false,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to set QoS: %w",
			err,
		)
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

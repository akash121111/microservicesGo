package messaging

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

const (
	RabbitMQURL = "amqp://guest:guest@localhost:5672/"

	OrderExchange = "order.events"

	// Commands / Events
	OrderCreatedRoutingKey           = "order.created"
	ReserveStockRoutingKey           = "stock.reserve"
	StockReservedRoutingKey          = "stock.reserved"
	ReservationStockFailedRoutingKey = "stock.reservation.failed"
	PaymentProcessRoutingKey         = "payment.process"
	PaymentSuccessRoutingKey         = "payment.success"
	PaymentFailedRoutingKey          = "payment.failed"

	// Separate Saga queues
	OrderSagaQueue   = "saga.order.queue"
	StockSagaQueue   = "saga.stock.queue"
	PaymentSagaQueue = "saga.payment.queue"

	SagaDLX = "saga.dlx"
	SagaDLQ = "saga.dlq"
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
	// Saga DLX
	// ========================================

	err = r.Ch.ExchangeDeclare(
		SagaDLX,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to declare saga DLX: %w",
			err,
		)
	}

	// ========================================
	// Saga DLQ
	// ========================================

	_, err = r.Ch.QueueDeclare(
		SagaDLQ,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to declare saga DLQ: %w",
			err,
		)
	}

	err = r.Ch.QueueBind(
		SagaDLQ,
		SagaDLQ,
		SagaDLX,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to bind saga DLQ: %w",
			err,
		)
	}

	// ========================================
	// Order Saga Queue
	// ========================================

	_, err = r.Ch.QueueDeclare(
		OrderSagaQueue,
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-dead-letter-exchange":    SagaDLX,
			"x-dead-letter-routing-key": SagaDLQ,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"failed to declare order saga queue: %w",
			err,
		)
	}

	// order.created -> OrderSagaQueue

	err = r.Ch.QueueBind(
		OrderSagaQueue,
		OrderCreatedRoutingKey,
		OrderExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to bind order saga queue: %w",
			err,
		)
	}

	// ========================================
	// Stock Saga Queue
	// ========================================

	_, err = r.Ch.QueueDeclare(
		StockSagaQueue,
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-dead-letter-exchange":    SagaDLX,
			"x-dead-letter-routing-key": SagaDLQ,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"failed to declare stock saga queue: %w",
			err,
		)
	}

	// stock.reserved -> StockSagaQueue

	err = r.Ch.QueueBind(
		StockSagaQueue,
		StockReservedRoutingKey,
		OrderExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to bind stock.reserved: %w",
			err,
		)
	}

	// stock.reservation.failed -> StockSagaQueue

	err = r.Ch.QueueBind(
		StockSagaQueue,
		ReservationStockFailedRoutingKey,
		OrderExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to bind stock.reservation.failed: %w",
			err,
		)
	}

	// ========================================
	// Payment Saga Queue
	// ========================================

	_, err = r.Ch.QueueDeclare(
		PaymentSagaQueue,
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-dead-letter-exchange":    SagaDLX,
			"x-dead-letter-routing-key": SagaDLQ,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"failed to declare payment saga queue: %w",
			err,
		)
	}

	// payment.success -> PaymentSagaQueue

	err = r.Ch.QueueBind(
		PaymentSagaQueue,
		PaymentSuccessRoutingKey,
		OrderExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to bind payment.success: %w",
			err,
		)
	}

	// payment.failed -> PaymentSagaQueue

	err = r.Ch.QueueBind(
		PaymentSagaQueue,
		PaymentFailedRoutingKey,
		OrderExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to bind payment.failed: %w",
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

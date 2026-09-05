package messaging

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

const (
	OrderExchange = "order.events"

	OrderCreatedRoutingKey = "order.created"

	NotificationQueue = "notification.queue"

	NotificationDLX = "notification.dlx"
	NotificationDLQ = "notification.dlq"
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
	return r.Ch.ExchangeDeclare(OrderExchange, "topic", true, false, false, false, nil)
}

func (r *RabbitMQ) Close() {
	if r.Ch != nil {
		_ = r.Ch.Close()
	}
	if r.Conn != nil {
		_ = r.Conn.Close()
	}

}

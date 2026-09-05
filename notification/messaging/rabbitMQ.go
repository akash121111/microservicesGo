package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	RabbitMQURL = "amqp://guest:guest@localhost:5672/"

	OrderExchange = "order.events"

	OrderCreatedRoutingKey = "order.created"

	NotificationQueue = "notification.queue"

	NotificationDLX = "notification.dlx"

	NotificationDLQ = "notification.dlq"
)

type RabbitMQConsumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQConsumer(
	url string,
) (*RabbitMQConsumer, error) {

	conn, err := amqp.Dial(url)

	if err != nil {
		return nil, fmt.Errorf(
			"rabbitmq connection failed: %w",
			err,
		)
	}

	ch, err := conn.Channel()

	if err != nil {
		conn.Close()

		return nil, fmt.Errorf(
			"rabbitmq channel creation failed: %w",
			err,
		)
	}

	r := &RabbitMQConsumer{
		conn: conn,
		ch:   ch,
	}

	if err := r.setup(); err != nil {
		r.Close()
		return nil, err
	}

	return r, nil
}

func (r *RabbitMQConsumer) setup() error {

	// --------------------------------
	// Main exchange
	// --------------------------------

	err := r.ch.ExchangeDeclare(
		OrderExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	// --------------------------------
	// Dead Letter Exchange
	// --------------------------------

	err = r.ch.ExchangeDeclare(
		NotificationDLX,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	// --------------------------------
	// Dead Letter Queue
	// --------------------------------

	_, err = r.ch.QueueDeclare(
		NotificationDLQ,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	// --------------------------------
	// Bind DLQ
	// --------------------------------

	err = r.ch.QueueBind(
		NotificationDLQ,
		NotificationDLQ,
		NotificationDLX,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	// --------------------------------
	// Main Queue
	// --------------------------------

	_, err = r.ch.QueueDeclare(
		NotificationQueue,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    NotificationDLX,
			"x-dead-letter-routing-key": NotificationDLQ,
		},
	)

	if err != nil {
		return err
	}

	// --------------------------------
	// Bind main queue
	// --------------------------------

	err = r.ch.QueueBind(
		NotificationQueue,
		OrderCreatedRoutingKey,
		OrderExchange,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	// --------------------------------
	// Prefetch
	// --------------------------------

	err = r.ch.Qos(
		10,
		0,
		false,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQConsumer) Consume() error {

	messages, err := r.ch.Consume(
		NotificationQueue,

		"",

		false, // manual ACK

		false,

		false,

		false,

		nil,
	)

	if err != nil {
		return fmt.Errorf(
			"consume failed: %w",
			err,
		)
	}

	log.Println(
		"Notification consumer started",
	)

	for message := range messages {

		var event OrderCreatedEvent

		err := json.Unmarshal(
			message.Body,
			&event,
		)

		if err != nil {

			log.Printf(
				"Invalid message: %v",
				err,
			)

			// Do NOT retry invalid JSON forever.
			//
			// false = don't requeue
			//
			// Because queue has DLX configured,
			// RabbitMQ sends it to DLQ.

			if err := message.Nack(
				false,
				false,
			); err != nil {
				log.Println(
					"failed to NACK:",
					err,
				)
			}

			continue
		}

		// --------------------------------
		// Process notification
		// --------------------------------

		// log.Printf(
		// 	"Notification received: orderId=%s userId=%s amount=%.2f",
		// 	event.OrderID,
		// 	event.UserID,
		// 	event.TotalAmount,
		// )

		log.Println(event)

		// --------------------------------
		// ACK
		// --------------------------------

		err = message.Ack(false)

		if err != nil {
			log.Println(
				"failed to ACK:",
				err,
			)
		}
	}

	return nil
}

func (r *RabbitMQConsumer) Close() {

	if r.ch != nil {
		_ = r.ch.Close()
	}

	if r.conn != nil {
		_ = r.conn.Close()
	}
}

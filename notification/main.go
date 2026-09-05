package main

import (
	"log"
	"notification/messaging"
)

func main() {

	rabbitMQ, err := messaging.NewRabbitMQConsumer(
		messaging.RabbitMQURL,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer rabbitMQ.Close()

	log.Println(
		"Notification Service started",
	)

	if err := rabbitMQ.Consume(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"payment/config"
	"payment/database"
	"payment/messaging"
	"payment/repository"
	"payment/service"

	"github.com/google/uuid"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("unable to load environment variable")
	}

	db, dbErr := database.NewDatabaseCongig(cfg)

	if dbErr != nil {
		log.Fatal("database is unable to connect")
	}
	migrationErr := database.MigrateDB(db)
	if migrationErr != nil {
		log.Fatal("database Migration is failed")
	}
	//event repository
	eventRepository := repository.NewProcessEventRepository(db)

	rabbitMQ, err := messaging.NewRabbitMQ(
		messaging.RabbitMQURL,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer rabbitMQ.Close()

	paymentService := service.NewPaymentService()

	// Temporary wallet for testing
	userID := uuid.MustParse(
		"8829fe9b-3ebc-4e53-8178-b845b4c1d9b6",
	)

	paymentService.CreateWallet(
		userID,
		100000,
	)

	log.Println("Payment Service started")

	if err := rabbitMQ.ConsumePayment(
		paymentService, eventRepository,
	); err != nil {

		log.Fatal(err)
	}
}

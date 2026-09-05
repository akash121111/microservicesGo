package main

import (
	"log"
	"net/http"
	"saga/client"
	"saga/config"
	"saga/database"
	"saga/messaging"
	"saga/repository"
	"saga/service"
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
	//repository for proces event

	eventRepository := repository.NewProcessEventRepository(db)

	rabbitMQ, err := messaging.NewRabbitMQ(
		messaging.RabbitMQURL,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitMQ.Close()

	orderClient := client.NewOrderClient(
		http.DefaultClient,
		"http://localhost:3003",
	)

	orderService := service.NewOrderSagaService(orderClient)

	productClient := client.NewProductClient(
		http.DefaultClient,
		"http://localhost:3002",
	)

	productService := service.NewProductSagaService(productClient)

	log.Println("Starting order.created consumer")

	go func() {
		if err := rabbitMQ.ConsumeOrderCreated(*eventRepository); err != nil {
			log.Println("order consumer error:", err)
		}
	}()

	log.Println("Starting stock event consumer")

	go func() {
		if err := rabbitMQ.ConsumeStockEvent(orderService, eventRepository); err != nil {
			log.Println("stock consumer error:", err)
		}
	}()

	go func() {
		if err := rabbitMQ.ConsumePaymentEvent(orderService, productService, eventRepository); err != nil {
			log.Println("stock consumer error:", err)
		}
	}()

	log.Println("Saga Service started")

	select {}
}

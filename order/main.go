package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"order/client"
	"order/config"
	"order/database"
	"order/handler"
	"order/logger"
	"order/messaging"
	"order/middleware"
	"order/redis"
	"order/repository"
	"order/routes"
	"order/service"
	"order/worker"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.Correlation())

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

	redisClient := redis.NewRedisClient(cfg)
	if redisClient == nil {
		log.Fatal("Redis error")
	}
	log.Println("Database connected Successfully")

	rabbitMQURL := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.RABBITMQ_USER,
		cfg.RABBITMQ_PASSWORD,
		cfg.RABBITMQ_HOST,
		cfg.RABBITMQ_PORT,
	)
	rabbitMQ, err := messaging.NewRabbitMQ(
		rabbitMQURL,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer rabbitMQ.Close()

	log.Println("RabbitMQ connected")

	//health Handler
	healthHandler := handler.NewHealthHandler(db, redisClient)
	router.GET("/order/health/live", healthHandler.Health)
	router.GET("/order/health/ready", healthHandler.Ready)

	logger := logger.New("order-service")

	//client
	productClient, err := client.NewProductClient("localhost:50051")
	if err != nil {
		log.Fatal(err)
	}

	//order
	orderRepository := repository.NewOrderRepository(db)
	outboxEventRepository := repository.NewOutboxEventRepository(db)
	orderService := service.NewOrderService(db, orderRepository, productClient, outboxEventRepository, redisClient, rabbitMQ, logger)
	orderHandler := handler.NewOrderHandler(orderService)
	outboxWorker := worker.NewOutboxWorker(
		outboxEventRepository,
		rabbitMQ,
	)

	routes.Routes(router, orderHandler)

	server := &http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server Failed %v", err)
		}
	}()
	ctx := context.Background()

	go outboxWorker.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shut down signal recieved")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shut down by forced %v", err)
	}

	log.Printf("server stopeed")

	// router.Run(":" + cfg.Port)
}

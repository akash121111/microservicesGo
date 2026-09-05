package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"product/config"
	"product/database"
	productgrpc "product/grpc"
	"product/handler"
	"product/logger"
	"product/messaging"
	"product/middleware"
	"product/redis"
	"product/repository"
	"product/routes"
	"product/service"

	productpb "proto/product"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc"
)

func main() {
	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.Correlation())
	//logger
	logger := logger.New("product-service")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("unable to load environment variable")
	}

	rabbitMQURL := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.RABBITMQ_USER,
		cfg.RABBITMQ_PASSWORD,
		cfg.RABBITMQ_HOST,
		cfg.RABBITMQ_PORT,
	)
	log.Println(rabbitMQURL)
	rabbitMQ, err := messaging.NewRabbitMQ(
		rabbitMQURL,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer rabbitMQ.Close()

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
	router.GET("/health", func(ctx *gin.Context) {
		time.Sleep(3 * time.Second)
		ctx.JSON(200, gin.H{
			"success": true,
			"message": "Health ok",
		})
	})
	//health Handler
	healthHandler := handler.NewHealthHandler(db, redisClient)
	router.GET("/product/health/live", healthHandler.Health)
	router.GET("/product/health/ready", healthHandler.Ready)
	//process event Repository
	processEvent := repository.NewProcessEventRepository(db)

	//product
	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository, redisClient, db, logger)
	grpcServer := grpc.NewServer()
	productGRPCServer := productgrpc.NewProductServer(productService)
	productpb.RegisterProductServiceServer(grpcServer, productGRPCServer)

	grpcListner, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen on gRPC port: %v", err)
	}

	go func() {
		log.Println("gRPC server listening on :", cfg.GRPCPort)

		if err := grpcServer.Serve(grpcListner); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	productHandler := handler.NewProductHandler(productService)

	routes.Routes(router, productHandler)

	server := &http.Server{
		Addr:    cfg.Port,
		Handler: router,
	}

	// Start RabbitMQ consumer in background
	go func() {
		if err := rabbitMQ.ConsumeReserveStock(productService, processEvent); err != nil {
			log.Printf("RabbitMQ consumer stopped: %v", err)
		}
	}()
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server Failed %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shut down signal recieved")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shut down by forced %v", err)
	}
	grpcServer.GracefulStop()
	log.Printf("server stopeed")

	// router.Run(":" + cfg.Port)
}

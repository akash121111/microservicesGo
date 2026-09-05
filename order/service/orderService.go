package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"order/apperror"
	"order/client"
	"order/messaging"
	"order/middleware"
	"order/model"
	"order/repository"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo     *repository.OrderRepository
	productClient *client.ProductClient
	eventRepo     *repository.OutboxEventRepository
	db            *gorm.DB
	redis         *redis.Client
	rabbitMQ      *messaging.RabbitMQ
	logger        *slog.Logger
}

func NewOrderService(db *gorm.DB, orderRepo *repository.OrderRepository, productClient *client.ProductClient, eventRepo *repository.OutboxEventRepository, redis *redis.Client, rebbitMQ *messaging.RabbitMQ, logger *slog.Logger) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		db:            db,
		productClient: productClient,
		eventRepo:     eventRepo,
		redis:         redis,
		rabbitMQ:      rebbitMQ,
		logger:        logger,
	}
}

func (s *OrderService) GetMyAllOrderWithItem(ctx context.Context, userID string) []model.OrderDataWithItems {

	return s.orderRepo.GetMyAllOrderWithItem(ctx, userID)
}

func (s *OrderService) GetOrderById(ctx context.Context, id string) (*model.OrderDataWithItems, error) {
	correlationId := middleware.GetCorrelationID(ctx)
	s.logger.Info(
		"fetch order",
		"orderId", id,
		"correlationID", correlationId,
	)
	order, err := s.orderRepo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, apperror.ErrOrderNotFound
	}
	// if order.Order.UserID != userId && role != "ADMIN" {
	// 	return nil, apperror.ErrForbidden
	// }
	return order, nil
}

func (s *OrderService) CreateOrder(
	ctx context.Context,
	idempotencyKey string,
	userID uuid.UUID,
	orderRequest model.OrderRequestBody,
) (*model.OrderModel, error) {

	correlationId := middleware.GetCorrelationID(ctx)
	s.logger.Info(
		"create order",
		"correlationID", correlationId,
	)
	items := make([]model.OrderItemModel, 0, len(orderRequest.Item))
	totalPrice := float64(0)
	if idempotencyKey == "" {
		return nil, apperror.ErrIdempotencyKey
	}
	key := fmt.Sprintf("ideempotency:order:%s:%s", userID.String(), idempotencyKey)
	cashe, err := s.redis.Get(ctx, key).Result()
	if err == nil {
		var orderResult model.OrderModel
		if cashe == "processing" {
			return nil, apperror.ErrOrderProcessing
		}
		if err := json.Unmarshal([]byte(cashe), &orderResult); err != nil {
			return nil, err
		}
		return &orderResult, nil
	}

	if err != redis.Nil {
		return nil, err
	}

	// 2. Atomically claim the idempotency key
	claim, err := s.redis.SetNX(ctx, key, "processing", 10*time.Minute).Result()
	if err != nil {
		return nil, err
	}
	if !claim {
		return nil, apperror.ErrOrderProcessing
	}
	// 1. Get products and validate stock
	for _, item := range orderRequest.Item {

		product, err := s.productClient.GetProduct(
			ctx,
			item.ProductID,
		)
		if err != nil {
			s.redis.Del(ctx, key)
			return nil, err
		}

		productID, err := uuid.Parse(product.ID)
		if err != nil {
			s.redis.Del(ctx, key)
			return nil, err
		}

		items = append(items, model.OrderItemModel{
			ProductID: productID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})

		totalPrice += product.Price * float64(item.Quantity)
	}

	// 2. Create Order + OrderItems in one transaction
	var result model.OrderModel

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		orderRepoTx := s.orderRepo.WithTx(tx)
		eventRepoTx := s.eventRepo.WithTx(tx)

		order := model.OrderModel{
			UserID:         userID,
			Status:         model.OrderPending,
			Total:          totalPrice,
			IdempotencyKey: idempotencyKey,
		}

		orderData, err := orderRepoTx.CreateOrder(ctx, order)
		if err != nil {
			return err
		}

		// Set OrderID for all items
		for i := range items {
			items[i].OrderID = orderData.ID
		}

		// IMPORTANT: use transaction repository
		_, err = orderRepoTx.CreateOrderItem(ctx, items)
		if err != nil {
			return err
		}

		result = *orderData

		orderEventItem := make([]messaging.OrderEventItem, 0, len(items))
		for _, itemsEvent := range items {
			orderEventItem = append(orderEventItem, messaging.OrderEventItem{
				ProductID: itemsEvent.ProductID,
				Quantity:  itemsEvent.Quantity,
			})
		}
		eventId := uuid.New()
		event := messaging.OrderCreatedEvent{
			EventID:     eventId,
			Event:       messaging.OrderCreatedRoutingKey,
			OrderID:     result.ID,
			UserID:      result.UserID,
			TotalAmount: totalPrice,
			Items:       orderEventItem,
		}
		log.Println(event)
		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}
		err = eventRepoTx.CreateOutBoxEvent(ctx, model.OutboxEvent{
			ID:         eventId,
			EventType:  "OrderCreated",
			RoutingKey: messaging.OrderCreatedRoutingKey,
			Payload:    payload,
			Published:  false,
		})
		if err != nil {
			return err

		}

		return nil
	})

	if err != nil {
		s.redis.Del(ctx, key)
		return nil, err
	}

	// //send message queue"
	// orderEventItem := make([]messaging.OrderEventItem, 0, len(items))
	// for _, itemsEvent := range items {
	// 	orderEventItem = append(orderEventItem, messaging.OrderEventItem{
	// 		ProductID: itemsEvent.ProductID,
	// 		Quantity:  itemsEvent.Quantity,
	// 	})
	// }
	// event := messaging.OrderCreatedEvent{
	// 	Event:       messaging.OrderCreatedRoutingKey,
	// 	OrderID:     result.ID,
	// 	UserID:      result.UserID,
	// 	TotalAmount: totalPrice,
	// 	Items:       orderEventItem,
	// }

	// err = s.rabbitMQ.PublishOrderCreate(ctx, event)

	// if err != nil {

	// 	// IMPORTANT:
	// 	// Order already exists.
	// 	// RabbitMQ publish failed.

	// 	return &result, fmt.Errorf(
	// 		"order created but event publishing failed: %w",
	// 		err,
	// 	)
	// }

	//store to redis
	dataCashe, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	err = s.redis.Set(ctx, key, dataCashe, 24*time.Hour).Err()
	if err != nil {
		log.Println("redis failed", err)
	}
	return &result, nil
}
func (s *OrderService) CancelOrder(
	ctx context.Context,
	orderID string,
	// userID uuid.UUID,
	// role string,
) error {

	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return apperror.ErrOrderNotFound
	}

	// Owner OR admin can cancel
	// if order.Order.UserID != userID && role != "ADMIN" {
	// 	return nil, apperror.ErrForbidden
	// }

	if order.Order.Status == model.OrderCancelled {
		return nil
	}

	// Release stock first
	// for _, orderItem := range order.Items {

	// 	err := s.productClient.ReleaseStock(
	// 		ctx,
	// 		orderItem.ProductID.String(),
	// 		int(orderItem.Quantity),
	// 	)

	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	// Change order status only after stock release succeeds
	_, err = s.orderRepo.UpdateOrderStatus(
		ctx,
		orderID,
		string(model.OrderCancelled),
	)

	if err != nil {
		// TODO:
		// compensate stock release here
		return err
	}

	return nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, id string, status string) error {
	_, err := s.orderRepo.UpdateOrderStatus(
		ctx,
		id,
		status,
	)
	if err != nil {
		return err
	}
	return nil
}

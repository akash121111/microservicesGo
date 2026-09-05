package service

import (
	"context"
	"saga/client"
	"saga/model"

	"github.com/google/uuid"
)

type OrderSagaService struct {
	orderClient *client.OrderClient
}

func NewOrderSagaService(orderClient *client.OrderClient) *OrderSagaService {
	return &OrderSagaService{
		orderClient: orderClient,
	}
}

func (s *OrderSagaService) CancelOrder(ctx context.Context, orderId uuid.UUID) error {
	return s.orderClient.CancelOrder(ctx, orderId)
}

func (s *OrderSagaService) UpdateOrderStatus(ctx context.Context, orderId uuid.UUID, status string) error {
	return s.orderClient.UpdateOrderStatus(ctx, orderId, status)
}

func (s *OrderSagaService) GetOrderByID(ctx context.Context, orderId uuid.UUID) (model.OrderDataWithItems, error) {
	return s.orderClient.GetOrderById(ctx, orderId)
}

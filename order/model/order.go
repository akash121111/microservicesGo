package model

import (
	"github.com/google/uuid"
)

type OrderModel struct {
	ID             uuid.UUID   `json:"id"`
	UserID         uuid.UUID   `json:"userID"`
	Total          float64     `json:"total" binding:"required"`
	Status         OrderStatus `json:"status"`
	IdempotencyKey string      `json:"idempotencyKey"`
}

type OrderItemModel struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"orderID"`
	ProductID uuid.UUID `json:"productID"`
	Quantity  int64     `json:"quantity" binding:"required"`
	Price     float64   `json:"price" binding:"required"`
}
type OrderRequest struct {
	ProductID string `json:"productID"`
	Quantity  int64  `json:"quantity" binding:"required"`
}

type OrderRequestBody struct {
	Item []OrderRequest `json:"item"`
}

type ProductModel struct {
	ID          string  `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Stock       int64   `json:"stock" binding:"required,gt=0"`
	Price       float64 `json:"price"  binding:"required,gt=0"`
	Status      string  `json:"status"`
}

type APIResponse[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}
type OrderPatchRequest struct {
	Status *string `json:"status"`
}

type OrderDataWithItems struct {
	Order OrderModel       `json:"order"`
	Items []OrderItemModel `json:"items"`
}

// type OutboxEventModel struct {
// 	ID         uuid.UUID `json:"id"`
// 	EventType  string    `json:"eventType"`
// 	RoutingKey string    `json:"routingKey"`
// 	Payload    []byte    `json:"payload"`
// 	Published  bool      `json:"published"`

// }

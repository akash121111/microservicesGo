package messaging

import (
	"product/model"

	"github.com/google/uuid"
)

type ReserveStockCommand struct {
	Event       string                 `json:"event"`
	EventID     uuid.UUID              `json:"eventId"`
	OrderID     uuid.UUID              `json:"orderId"`
	UserID      uuid.UUID              `json:"userId"`
	Items       []model.OrderEventItem `json:"items"`
	TotalAmount float64                `json:"totalAmount"`
}

type ReserveStockFailedEvent struct {
	Event   string                 `json:"event"`
	EventID uuid.UUID              `json:"eventId"`
	OrderID uuid.UUID              `json:"orderId"`
	Items   []model.OrderEventItem `json:"items"`
	Reason  string                 `json:"reason"`
}

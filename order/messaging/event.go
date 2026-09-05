package messaging

import "github.com/google/uuid"

type OrderCreatedEvent struct {
	EventID     uuid.UUID        `json:"eventID"`
	Event       string           `json:"event"`
	OrderID     uuid.UUID        `json:"orderId"`
	UserID      uuid.UUID        `json:"userId"`
	TotalAmount float64          `json:"totalAmount"`
	Items       []OrderEventItem `json:"items"`
}
type OrderEventItem struct {
	ProductID uuid.UUID `json:"productID"`
	Quantity  int64     `json:"quantity"`
}
